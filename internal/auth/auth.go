package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"main/internal/config"
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

func GenerateTokens(userId string, refreshTokenCreator RefreshTokenCreator, cfg *config.Config, tokenAuth *jwtauth.JWTAuth) (Tokens, error) {
	const op = "auth.GenerateTokens"

	id := uuid.New().String()

	accessExp := time.Now().Add(cfg.Authorization.AccessTTL)
	refreshExp := time.Now().Add(cfg.Authorization.RefreshTTL)

	_, refreshToken, err := tokenAuth.Encode(map[string]interface{}{
		"token_id": id,
		"exp":      refreshExp,
	})
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	hashedRefreshToken := HashRefreshToken(refreshToken, cfg.Authorization.Salt)

	err = refreshTokenCreator.CreateRefreshToken(id, userId, hashedRefreshToken, refreshExp)
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	_, accessToken, err := tokenAuth.Encode(map[string]interface{}{
		"user_id": userId,
		"exp":     accessExp,
	})
	if err != nil {
		return Tokens{}, fmt.Errorf("%s: %w", op, err)
	}

	return Tokens{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}, nil
}

func HashRefreshToken(token, salt string) string {
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}
