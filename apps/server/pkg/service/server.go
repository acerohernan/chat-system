package service

import (
	"net/http"
	"time"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/logger"
	"github.com/chat-system/server/pkg/rtc"
)

type ChatServer struct {
	rtcService *rtc.RTCService
	httpServer *http.Server
}

func NewChatServer(config *config.Config, rtcService *rtc.RTCService) *ChatServer {
	return &ChatServer{
		rtcService: rtcService,
	}
}

func (c *ChatServer) Start() error {
	// real time communication
	http.HandleFunc("/rtc", c.rtcService.ServeHTTP)

	server := &http.Server{
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()

	if err != nil {
		logger.Errorw("couldn't start http server", err)
		return err
	}

	return nil
}

func (*ChatServer) Stop(kill bool) error {
	return nil
}
