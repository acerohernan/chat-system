package service

import (
	"sync"

	core "github.com/chat-system/server/proto"
)

type LocalStorage struct {
	publickKeys map[core.UserEmail]*core.PublicKey
	mu          sync.RWMutex
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		publickKeys: make(map[core.UserEmail]*core.PublicKey),
		mu:          sync.RWMutex{},
	}
}

func (s *LocalStorage) StorePublicKey(key *core.PublicKey) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.publickKeys[core.UserEmail(key.UserEmail)] = key
	return nil
}

func (s *LocalStorage) GetPublicKey(email core.UserEmail) (*core.PublicKey, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	key := s.publickKeys[email]

	if key == nil {
		return nil, ErrPublicKeyNotFound
	}

	return key, nil
}
