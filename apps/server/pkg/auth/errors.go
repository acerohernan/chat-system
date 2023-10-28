package auth

import "errors"

var (
	ErrInvalidAccessToken = errors.New("invalid access token")
)