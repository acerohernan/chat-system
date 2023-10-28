package service

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/logger"
	"github.com/chat-system/server/pkg/rtc"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type ChatServer struct {
	config      *config.Config
	rtcService  *rtc.RTCService
	authService *AuthService
	httpServer  *http.Server
	storage     PersistentStorage
	running     atomic.Bool
	doneChann   chan struct{}
	closeChann  chan struct{}
}

func NewChatServer(config *config.Config) (*ChatServer, error) {
	s := &ChatServer{
		config:     config,
		running:    atomic.Bool{},
		doneChann:  make(chan struct{}),
		closeChann: make(chan struct{}),
	}

	s.rtcService = rtc.NewRTCService()

	mc, err := getMongoClient(config.Mongo)

	if err != nil {
		return nil, err
	}

	s.storage = NewMongoStorage(config.Mongo, mc)

	if err != nil {
		return nil, err
	}

	s.authService = NewAuthService(config, s.storage)

	middlewares := []negroni.Handler{
		// always first
		negroni.NewRecovery(),
		// CORS is allowed, we rely on token authentication to prevent improper use
		cors.New(cors.Options{
			AllowOriginFunc: func(origin string) bool {
				return true
			},
			AllowedHeaders: []string{"*"},
			// allow preflight to be cached for a day
			MaxAge: 86400,
		}),
	}

	mux := mux.NewRouter()

	// health check
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("OK")) })

	// auth
	mux.HandleFunc("/auth/complete", s.authService.CompleteRegistrationHTTP).Methods("POST")
	mux.HandleFunc("/auth/{provider}", s.authService.BeginAuthHTTP)
	mux.HandleFunc("/auth/{provider}/callback", s.authService.AuthCallbackHTTP)

	// websocket
	mux.HandleFunc("/rtc", s.rtcService.ServeHTTP)

	s.httpServer = &http.Server{
		Handler: configureMiddlewares(mux, middlewares...),
	}

	return s, nil
}

func (c *ChatServer) Start() error {
	logger.Infow("starting server...")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", c.config.Port))

	if err != nil {
		return err
	}

	go c.httpServer.Serve(listener)

	logger.Infow("http server running", "port", c.config.Port)

	c.running.Store(true)

	logger.Infow("server started successfully")

	<-c.doneChann

	// shutdown
	_ = c.httpServer.Shutdown(context.Background())

	_ = c.storage.Close()

	close(c.closeChann)

	return nil
}

func (c *ChatServer) Stop(force bool) {
	if !c.running.Swap(false) {
		return
	}

	close(c.doneChann)

	<-c.closeChann

	return
}

func configureMiddlewares(handler http.Handler, middlewares ...negroni.Handler) *negroni.Negroni {
	n := negroni.New()
	for _, m := range middlewares {
		n.Use(m)
	}
	n.UseHandler(handler)
	return n
}
