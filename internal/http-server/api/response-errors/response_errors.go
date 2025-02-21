package responseerrors

import (
	"errors"
)

var (
	ErrUserNotOwner        = errors.New("user is not owner or note is not exists")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrInternalServerError = errors.New("internal server error")

	ErrUserUnauthorized = errors.New("user unauthorized")

	ErrAccessTokenExpired  = errors.New("access token expired")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	ErrUserDoesNotExist    = errors.New("user does not exist")
	ErrUserIsAlreadyExists = errors.New("user is already exists")
	ErrInvalidPassword     = errors.New("invalid password")

	ErrFailedToAddNoteNode = errors.New("failed to add note node")
	ErrFailedToDeleteNode  = errors.New("failed to delete node")
	ErrNodeDoesNotExist    = errors.New("node does not exist")
)
