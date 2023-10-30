package service

import (
	"context"

	core "github.com/chat-system/server/proto"
)

type PersistentStorage interface {
	StoreUser(user *core.User) error
	GetUser(id core.UserId) (*core.User, error)
	GetUserWithEmail(email core.UserEmail) (*core.User, error)

	// close active client connections
	Close(ctx context.Context) error
}

type InMemoryStorage interface {
}
