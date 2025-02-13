package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/auth"
	"main/internal/storage"
	"time"
)

func (s *Storage) CreateRefreshToken(id, userId, tokenHash string, expiresAt time.Time) error {
	const op = "storage.postgres.CreateRefreshToken"

	_, err := s.db.Exec(createRefreshTokenQuery, id, userId, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetUserByTokenId(id string) (auth.RefreshToken, error) {
	const op = "storage.postgres.GetRefreshToken"

	var refreshToken auth.RefreshToken

	err := s.db.Get(&refreshToken, getRefreshTokenByIdQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.RefreshToken{}, storage.ErrRefreshTokenNotFound
	}
	if err != nil {
		return auth.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}

func (s *Storage) RevokeRefreshTokenById(id string) error {
	const op = "storage.postgres.RevokeRefreshTokenById"

	// revoking refresh token
	res, err := s.db.Exec(revokeRefreshTokenByIdQuery, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	// check if refresh token wasn't found
	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rows == 0 {
		return storage.ErrRefreshTokenNotFound
	}

	return nil
}

func (s *Storage) RevokeExpiredRefreshTokens() error {
	const op = "storage.postgres.RevokeExpiredRefreshTokens"

	// revoking expired refresh tokens
	_, err := s.db.Exec(revokeExpiredRefreshTokensQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteWeekOldRefreshTokens() error {
	const op = "storage.postgres.DeleteWeekOldRefreshTokens"

	// deleting week old refresh tokens
	_, err := s.db.Exec(deleteWeekOldRefreshTokensQuery)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
