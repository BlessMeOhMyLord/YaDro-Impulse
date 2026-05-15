package internal

import "errors"

var (
	ErrNotFound      = errors.New("dns server not exists")
	ErrIsIncorrect   = errors.New("dns server is incorrect")
	ErrAlreadyExists = errors.New("dns server is already exists")
)
