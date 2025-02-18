package responseerrors

import (
	"errors"
)

var (
	ErrUserNotOwner        = errors.New("user is not owner or note is not exists")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInternalServerError = errors.New("internal server error")
)
