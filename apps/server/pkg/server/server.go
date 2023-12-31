package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/config/logger"
	"github.com/chat-system/server/pkg/controllers"
	"github.com/chat-system/server/pkg/service"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type ChatServer struct {
	config         *config.Config
	rtcController  *controllers.RTCController
	authController *controllers.AuthController
	userController *controllers.UserController
	httpServer     *http.Server
	storage        service.PersistentStorage
	running        atomic.Bool
	doneChann      chan struct{}
	closeChann     chan struct{}
}

func NewChatServer(config *config.Config, rtcController *controllers.RTCController, authController *controllers.AuthController, userController *controllers.UserController, storage service.PersistentStorage) (*ChatServer, error) {
	s := &ChatServer{
		config:         config,
		running:        atomic.Bool{},
		doneChann:      make(chan struct{}),
		closeChann:     make(chan struct{}),
		storage:        storage,
		rtcController:  rtcController,
		authController: authController,
		userController: userController,
	}

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
	mux.HandleFunc("/auth/complete", s.authController.CompleteRegistration).Methods("POST")
	mux.HandleFunc("/auth/{provider}", s.authController.BeginAuth)
	mux.HandleFunc("/auth/{provider}/callback", s.authController.AuthCallback)

	// user
	mux.HandleFunc("/user", s.userController.FindUser).Methods("GET")

	// websocket
	mux.HandleFunc("/rtc", s.rtcController.Serve)

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
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	_ = c.httpServer.Shutdown(ctx)

	_ = c.storage.Close(ctx)

	cancel()

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
