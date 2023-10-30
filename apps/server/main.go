package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/config/logger"
	"github.com/chat-system/server/pkg/server"
)

func main() {
	config := config.NewConfig()

	logger.Init(config.Logger)

	server, err := server.InitializeServer(config)

	if err != nil {
		logger.Errorw("couldn't start the chat server", err)
		os.Exit(1)
	}

	sigChann := make(chan os.Signal, 1)
	signal.Notify(sigChann, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-sigChann
		logger.Infow("exit requested, shutting down...", "signal", sig)
		server.Stop(false)
	}()

	server.Start()
}
