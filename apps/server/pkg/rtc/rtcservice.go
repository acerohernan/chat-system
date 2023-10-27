package rtc

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chat-system/server/pkg/logger"
	protocol "github.com/chat-system/server/proto"
	"github.com/gorilla/websocket"
)

type RTCService struct {
	upgrader    websocket.Upgrader
	mu          sync.Mutex
	connections map[*websocket.Conn]struct{}
}

func NewRTCService() *RTCService {
	return &RTCService{
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// allow all origins, we use token for auth
				return true
			},
		},
		connections: make(map[*websocket.Conn]struct{}),
		mu:          sync.Mutex{},
	}
}

func (s *RTCService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// reject non websocket requests
	if !websocket.IsWebSocketUpgrade(r) {
		w.WriteHeader(404)
		return
	}

	// validate auth
	err := s.validateToken(r)

	if err != nil {
		// unauthorized
		s.handleError(w, http.StatusUnauthorized, err)
		return
	}

	// upgrade only once the basics are good to go
	conn, err := s.upgrader.Upgrade(w, r, nil)

	if err != nil {
		s.handleError(w, http.StatusInternalServerError, err)
		return
	}

	done := make(chan struct{})

	s.mu.Lock()
	s.connections[conn] = struct{}{}
	s.mu.Unlock()

	// function exits when websocket terminates, it'll close the event reading
	defer func() {
		logger.Infow("finishing WS connection...")
		s.mu.Lock()
		delete(s.connections, conn)
		s.mu.Unlock()
		close(done)
		logger.Infow("WS connection finished")
	}()

	wsClient := NewWsClient(conn)

	s.mu.Lock()
	logger.Infow("new client ws connected", "address", conn.RemoteAddr())
	s.mu.Unlock()

	// handle incoming request from websockets
	for {
		req, _, err := wsClient.ReadRequest()

		if err != nil {
			// normal/expected clousure
			if err == io.EOF ||
				strings.HasSuffix(err.Error(), "use of closed network connection") ||
				strings.HasSuffix(err.Error(), "connection reset by peer") ||
				websocket.IsCloseError(
					err,
					websocket.CloseAbnormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNormalClosure,
					websocket.CloseNoStatusReceived,
				) {
				logger.Infow("exit ws read loop for closed connection", "wsError", err)
			} else {
				logger.Errorw("error reading from websocket", err)
			}
			return
		}

		switch m := req.Message.(type) {
		case *protocol.SignalRequest_Ping:
			logger.Debugw("a ping request received", "message", m)
			_, _ = wsClient.WriteResponse(&protocol.SignalResponse{Message: &protocol.SignalResponse_Pong{
				Pong: &protocol.Pong{
					Timestamp: time.Now().UnixMilli(),
				},
			}})
		}
	}
}

func (s *RTCService) validateToken(r *http.Request) error {
	return nil
}

func (*RTCService) handleError(w http.ResponseWriter, status int, err error, keysAndValues ...interface{}) {
	logger.Errorw("error in handling connection", err)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}
