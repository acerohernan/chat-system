package controllers

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/chat-system/server/pkg/config"
	"github.com/chat-system/server/pkg/service"
	"github.com/chat-system/server/pkg/service/auth"
	"github.com/chat-system/server/pkg/utils"
	core "github.com/chat-system/server/proto"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

type AuthController struct {
	config   *config.AuthConfig
	storage  service.PersistentStorage
	verifier *auth.Verifier
}

func NewAuthController(config *config.Config, verifier *auth.Verifier, storage service.PersistentStorage) *AuthController {
	// setup oauth providers
	goth.UseProviders(
		google.New(config.Auth.GoogleClientId, config.Auth.GoogleClientSecret, config.Host+"/auth/google/callback", "email", "profile"),
	)

	store := sessions.NewCookieStore([]byte(config.Auth.JWTSecret))
	gothic.Store = store

	return &AuthController{
		config:   config.Auth,
		storage:  storage,
		verifier: verifier,
	}
}

func (s *AuthController) BeginAuth(w http.ResponseWriter, r *http.Request) {
	gothic.BeginAuthHandler(w, r)
	return
}

func (s *AuthController) AuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if user.Email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var UserId string
	var canSendMessage bool

	dbUser, err := s.storage.GetUserWithEmail(core.UserEmail(user.Email))

	if err != nil && err != service.ErrUserNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dbUser != nil {
		UserId = dbUser.Id
		canSendMessage = true
	} else {
		UserId = utils.NewGuid(utils.UserPrefix)
		// don't allow to send message until finish registration
		canSendMessage = false
	}

	token, err := s.verifier.CreateToken(&auth.Grants{
		Id:             UserId,
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

func (s *AuthController) CompleteRegistration(w http.ResponseWriter, r *http.Request) {
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

	// validate if param is a valid publick key
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

	dbUser, err := s.storage.GetUser(core.UserId(grants.Id))

	if err != nil && err != service.ErrUserNotFound {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if dbUser != nil {
		// if the user is created, reject request
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := &core.User{
		Id:        grants.Id,
		Email:     grants.Email,
		PublicKey: params.PublicKey,
	}

	err = s.storage.StoreUser(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *AuthController) validateToken(r *http.Request) (*auth.Grants, error) {
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
