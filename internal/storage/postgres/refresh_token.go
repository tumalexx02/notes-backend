package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"main/internal/auth"
	"main/internal/storage"
	"main/internal/storage/postgres/queries"
	"time"
)

func (s *Storage) CreateRefreshToken(id, userId, tokenHash string, expiresAt time.Time) error {
	const op = "storage.postgres.CreateRefreshToken"

	_, err := s.db.Exec(queries.CreateRefreshTokenQuery, id, userId, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetRefreshTokenById(id string) (auth.RefreshToken, error) {
	const op = "storage.postgres.GetRefreshTokenById"

	var refreshToken auth.RefreshToken

	err := s.db.Get(&refreshToken, queries.GetRefreshTokenByIdQuery, id)
	if errors.Is(err, sql.ErrNoRows) {
		return auth.RefreshToken{}, storage.ErrRefreshTokenNotFound
	}
	if err != nil {
		return auth.RefreshToken{}, fmt.Errorf("%s: %w", op, err)
	}

	return refreshToken, nil
}

func (s *Storage) DeleteExpiredRefreshTokens() (int, error) {
	const op = "storage.postgres.DeleteExpiredRefreshTokens"

	// deleting expired refresh tokens
	rows, err := s.db.Exec(queries.DeleteExpiredRefreshTokensQuery)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// counting rows
	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(rowsAffected), nil
}
