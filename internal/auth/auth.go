package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

type RefreshToken struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id" db:"user_id"`
	TokenHash string    `json:"token_hash" db:"token_hash"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	Revoked   bool      `json:"revoked"`
}

type Tokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

type RefreshTokenCreator interface {
	CreateRefreshToken(id, userId, tokenHash string, expiresAt time.Time) error
}

func GenerateTokens(refreshTokenCreator RefreshTokenCreator, userId string, refreshExpiresAt time.Time, accessExpiresAt time.Time, tokenAuth *jwtauth.JWTAuth) (Tokens, error) {
	const op = "auth.GenerateTokens"

	id := uuid.New().String()

	_, refreshToken, err := tokenAuth.Encode(map[string]interface{}{
		"token_id": id,
		"exp":      refreshExpiresAt,
	})
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	// TODO: add salt
	hashedRefreshToken := sha256.Sum256([]byte(refreshToken))

	err = refreshTokenCreator.CreateRefreshToken(id, userId, hex.EncodeToString(hashedRefreshToken[:]), refreshExpiresAt)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	_, accessToken, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": userId,
		"exp":     accessExpiresAt,
	})
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}
