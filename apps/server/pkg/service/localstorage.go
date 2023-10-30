package service

import (
	"context"
	"sync"

	core "github.com/chat-system/server/proto"
)

type LocalStorage struct {
	users        map[core.UserId]*core.User
	usersByEmail map[core.UserEmail]*core.User
	mu           sync.RWMutex
}

func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		users:        make(map[core.UserId]*core.User),
		usersByEmail: make(map[core.UserEmail]*core.User),
		mu:           sync.RWMutex{},
	}
}

func (s *LocalStorage) StoreUser(user *core.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.users[core.UserId(user.Id)] = user
	s.usersByEmail[core.UserEmail(user.Email)] = user
	return nil
}

func (s *LocalStorage) GetUser(id core.UserId) (*core.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user := s.users[id]

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *LocalStorage) GetUserWithEmail(email core.UserEmail) (*core.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user := s.usersByEmail[email]

	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (s *LocalStorage) Close(_ context.Context) error {
	return nil
}
