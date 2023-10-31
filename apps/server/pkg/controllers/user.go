package controllers

import (
	"encoding/json"
	"net/http"
	"net/mail"

	"github.com/chat-system/server/pkg/service"
	"github.com/chat-system/server/pkg/service/auth"
	core "github.com/chat-system/server/proto"
)

type UserController struct {
	storage  service.PersistentStorage
	verifier *auth.Verifier
}

func NewUserController(storage service.PersistentStorage, verifier *auth.Verifier) *UserController {
	return &UserController{
		storage:  storage,
		verifier: verifier,
	}
}

type FindUserParams struct {
	Email string `json:"email"`
}

func (c *UserController) FindUser(w http.ResponseWriter, r *http.Request) {
	token := auth.ExtractTokenFromRequest(r)

	_, err := c.verifier.ParseToken(token)

	if err != nil {
		w.WriteHeader(401)
		return
	}

	query := r.URL.Query()

	params := &FindUserParams{
		Email: query.Get("email"),
	}

	if err != nil || params.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// validate if it's a valid email
	_, err = mail.ParseAddress(params.Email)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := c.storage.GetUserWithEmail(core.UserEmail(params.Email))

	if err == service.ErrUserNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userJson, err := json.Marshal(user)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(userJson)
}
