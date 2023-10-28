package auth

import (
	"fmt"
	"net/http"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/logger"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct {
	config *config.AuthConfig
}

func NewAuthService(config *config.Config) *AuthService {
	// setup oauth providers
	goth.UseProviders(
		google.New(config.Auth.GoogleClientId, config.Auth.GoogleClientSecret, config.Host+"/auth/google/callback", "email", "profile"),
	)

	// do not store session in cookies, we'll handle sessions with jwt
	store := sessions.NewCookieStore([]byte(config.Auth.JWTSecret))
	gothic.Store = store

	return &AuthService{
		config: config.Auth,
	}
}

func (s *AuthService) BeginAuthHTTP(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
	return
}

func (s *AuthService) AuthCallbackHTTP(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user.Email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token := NewAccessToken(s.config.JWTIssuer, s.config.JWTSecret, &Grants{
		Email:          user.Email,
		CanSendMessage: true,
	})

	jwt, err := token.toJWT()

	if err != nil {
		logger.Errorw("error at creating jwt", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"jwt": %v}`, jwt)))
}
