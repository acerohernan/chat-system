package service

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/chat-system/server/pkg/auth"
	"github.com/chat-system/server/pkg/config"
	core "github.com/chat-system/server/proto"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthService struct {
	config            *config.AuthConfig
	persistentStorage PersistentStorage
	verifier          *auth.Verifier
}

func NewAuthService(config *config.Config, persistentStorage PersistentStorage) *AuthService {
	// setup oauth providers
	goth.UseProviders(
		google.New(config.Auth.GoogleClientId, config.Auth.GoogleClientSecret, config.Host+"/auth/google/callback", "email", "profile"),
	)

	// do not store session in cookies, we'll handle sessions with jwt
	store := sessions.NewCookieStore([]byte(config.Auth.JWTSecret))
	gothic.Store = store

	return &AuthService{
		config:            config.Auth,
		persistentStorage: persistentStorage,
		verifier:          auth.NewVerifier(config.Auth.JWTIssuer, config.Auth.JWTSecret),
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

	canSendMessage := true

	_, err = s.persistentStorage.GetPublicKey(core.UserEmail(user.Email))

	if err != nil {
		if err == ErrPublicKeyNotFound {
			canSendMessage = false
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	token, err := s.verifier.CreateToken(&auth.Grants{
		Email:          user.Email,
		CanSendMessage: canSendMessage,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"jwt": "%v"}`, token.Jwt)))
}

type CompleteRegistrationParams struct {
	PublicKey string `json:"publicKey"`
}

func (s *AuthService) CompleteRegistrationHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	grants, err := s.validateToken(r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var params CompleteRegistrationParams
	err = json.NewDecoder(r.Body).Decode(&params)

	if err != nil || params.PublicKey == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// validate that is a valid public key
	block, _ := pem.Decode([]byte(params.PublicKey))

	if block == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	key, err := s.persistentStorage.GetPublicKey(core.UserEmail(grants.Email))

	if err != nil && err != ErrPublicKeyNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if key != "" {
		// don't allow save a new public key to users that has already saved it
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.persistentStorage.StorePublicKey(core.UserEmail(grants.Email), core.UserPublicKey(params.PublicKey))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *AuthService) validateToken(r *http.Request) (*auth.Grants, error) {
	header := r.Header.Get("Authorization")

	if header == "" {
		return nil, auth.ErrInvalidAccessToken
	}

	rawJWT, prefixFound := strings.CutPrefix(header, "Bearer ")

	if !prefixFound {
		return nil, auth.ErrInvalidAccessToken
	}

	// validate authenticity
	accessToken, err := s.verifier.ParseToken(rawJWT)

	if err != nil {
		return nil, auth.ErrInvalidAccessToken
	}

	return accessToken.Grants, nil
}
