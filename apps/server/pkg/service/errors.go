package service

import "errors"

var (
	ErrPublicKeyNotFound = errors.New("public key for user not found in persistent storage")
)
