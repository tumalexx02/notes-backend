package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/models/user"
	"main/internal/storage"

	"github.com/lib/pq"
)

const (
	ErrUniqueViolation = "23505"
)

func (s *Storage) CreateUser(email, name, passwordHash string) (string, error) {
	const op = "storage.postgres.CreateUser"

	// creating user
	var id string

	err := s.db.Get(&id, createUserQuery, email, name, passwordHash)
	if err != nil {
		// check if user already exists
		var sqlxerr *pq.Error
		if errors.As(err, &sqlxerr) && sqlxerr.Code == ErrUniqueViolation {
			return "", storage.ErrUserAlreadyExists
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) GetUser(email string) (user.User, error) {
	const op = "storage.postgres.GetUser"

	// getting user by email
	var userFromDB user.User

	err := s.db.Get(&userFromDB, getUserByEmailQuery, email)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userFromDB, nil
}

func (s *Storage) GetUserById(id string) (user.User, error) {
	const op = "storage.postgres.GetUser"

	// getting user by id
	var userFromDB user.User

	err := s.db.Get(&userFromDB, getUserByIdQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return user.User{}, storage.ErrUserNotFound
	}
	if err != nil {
		return user.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return userFromDB, nil
}
