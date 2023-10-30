package service

import "errors"

var (
	ErrUserNotFound = errors.New("user not found in persistent storage")
)
