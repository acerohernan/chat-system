package service

import (
	"sync"

	core "github.com/chat-system/server/proto"
)

type LocalStorage struct {
	publickKeys map[core.UserEmail]core.UserPublicKey
	mu          sync.RWMutex
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		publickKeys: make(map[core.UserEmail]core.UserPublicKey),
		mu:          sync.RWMutex{},
	}
}

func (s *LocalStorage) StorePublicKey(email core.UserEmail, key core.UserPublicKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publickKeys[email] = key
	return nil
}

func (s *LocalStorage) GetPublicKey(email core.UserEmail) (core.UserPublicKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.publickKeys[email]

	if key == "" {
		return "", ErrPublicKeyNotFound
	}

	return key, nil
}
