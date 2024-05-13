package db

import "errors"

var (
	ErrStreamServerUnregistered = errors.New("stream server unregistered")
	ErrUserNotFound             = errors.New("user not found")
)
