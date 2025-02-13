package storage

import "errors"

var (
	ErrNoteNotFound     = errors.New("note not found")
	ErrNoteNodeNotFound = errors.New("note node not found")

	ErrUserAlreadyExists = errors.New("user with this email already exists")
	ErrUserNotFound      = errors.New("user not found")
)
