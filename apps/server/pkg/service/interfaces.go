package service

import core "github.com/chat-system/server/proto"

type PersistentStorage interface {
	StorePublicKey(email core.UserEmail, key core.UserPublicKey) error
	GetPublicKey(email core.UserEmail) (core.UserPublicKey, error)
}

type InMemoryStorage interface {
}
