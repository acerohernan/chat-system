package service

import core "github.com/chat-system/server/proto"

type PersistentStorage interface {
	StorePublicKey(key *core.PublicKey) error
	GetPublicKey(email core.UserEmail) (*core.PublicKey, error)

	// close active client connections
	Close() error
}

type InMemoryStorage interface {
}
