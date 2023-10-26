package rtc

import (
	"net/http"

	"github.com/chat-system/server/pkg/logger"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type RTCService struct {
}

func NewRTCService() *RTCService {
	return &RTCService{}
}

func (s *RTCService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		s.handleError(w, 400, err)
		return
	}

	// validate request (auth)

}

func (*RTCService) handleError(w http.ResponseWriter, status int, err error, keysAndValues ...interface{}) {
	logger.Errorw("error in handling connection", err)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}
