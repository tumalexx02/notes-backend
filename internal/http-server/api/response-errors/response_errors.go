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

	ErrFailedToAddNoteNode       = errors.New("failed to add note node")
	ErrFailedToDeleteNode        = errors.New("failed to delete node")
	ErrFailedToUpdateNodeContent = errors.New("failed to update node content")
	ErrNodeDoesNotExist          = errors.New("node does not exist")

	ErrNoteDoesNotExist         = errors.New("note does not exist")
	ErrNoteIsAlreadyExists      = errors.New("note is already exists")
	ErrFailedToArchiveNote      = errors.New("failed to archive note")
	ErrFailedToCreateNote       = errors.New("failed to create note")
	ErrFailedToDeleteNote       = errors.New("failed to delete note")
	ErrFailedToGetNote          = errors.New("failed to get note")
	ErrFailedToGetUsersNote     = errors.New("failed to get users note")
	ErrFailedToUnarchiveNote    = errors.New("failed to unarchive note")
	ErrFailedToUpdateFullNote   = errors.New("failed to update full note")
	ErrFailedToUpdateNodesOrder = errors.New("failed to update nodes order")
	ErrFailedToUpdateNoteTitle  = errors.New("failed to update note title")
)
