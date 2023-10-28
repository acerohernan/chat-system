package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chat-system/server/pkg/auth"
	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/logger"
	"github.com/chat-system/server/pkg/rtc"
	"github.com/chat-system/server/pkg/service"
)

func main() {
	config := config.NewConfig()

	logger.Init(config.Logger)

	rtcService := rtc.NewRTCService()

	authService := auth.NewAuthService(config)

	server := service.NewChatServer(config, rtcService, authService)

	sigChann := make(chan os.Signal, 1)
	signal.Notify(sigChann, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChann
		logger.Infow("exit requested, shutting down...", "signal", sig)
		server.Stop(false)
	}()

	server.Start()
}
