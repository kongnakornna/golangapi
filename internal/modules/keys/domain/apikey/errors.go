package apikey

import "errors"

var (
	ErrNotFound  = errors.New("api key not found")
	ErrRevoked   = errors.New("api key has been revoked")
	ErrExpired   = errors.New("api key has expired")
	ErrInvalid   = errors.New("api key is invalid")
	ErrForbidden = errors.New("access denied to this api key")
)
