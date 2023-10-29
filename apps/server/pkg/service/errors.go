package service

import "errors"

var (
	ErrPublicKeyNotFound = errors.New("public key for user not found in persistent storage")
	ErrUserNotFound      = errors.New("user not found in persistent storage")
)
