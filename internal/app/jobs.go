package app

import (
	"context"
	"log/slog"
	"time"
)

type TokenRevoker interface {
	DeleteExpiredRefreshTokens() (int, error)
}

func startTokensRevokingJob(ctx context.Context, log *slog.Logger, tokenRevoker TokenRevoker) {
	log.Info("token revoking job started")

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Info("token revoking job stopped", "reason", ctx.Err())
				return
			case <-ticker.C:
				revoked, err := tokenRevoker.DeleteExpiredRefreshTokens()
				if err != nil {
					log.Error("failed to revoke expired refresh tokens", "error", err)
					continue
				}

				log.Info("revoked expired refresh tokens", "count", revoked)
			}
		}
	}()
}
