package refresh

import (
	"log/slog"
	"main/internal/auth"
	"main/internal/config"
	resp "main/internal/http-server/api/response"
	"main/internal/http-server/api/validate"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/go-chi/render"
)

type Request struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type Response struct {
	resp.Response
	AccessToken string `json:"access_token"`
}

type RefreshTokener interface {
	GetRefreshTokenById(id string) (auth.RefreshToken, error)
	RevokeRefreshTokenById(id string) error
}

func New(cfg *config.Config, log *slog.Logger, refreshTokener RefreshTokener, tokenAuth *jwtauth.JWTAuth) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.auth.refresh.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request-id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := validate.DecodeAndValidateRequestJson(&req, w, r, log); err != nil {
			return
		}

		token, err := tokenAuth.Decode(req.RefreshToken)
		if err != nil {
			log.Error("failed to decode token", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to decode token"))

			return
		}

		tokenId, ok := token.Get("token_id")
		if !ok {
			log.Error("failed to get token id", slog.Attr{Key: "error", Value: slog.StringValue("failed to get token id")})

			render.JSON(w, r, resp.Error("failed to get token id"))

			return
		}

		tokenIdStr, ok := tokenId.(string)
		if !ok {
			log.Error("failed to convert token id into string", slog.Attr{Key: "error", Value: slog.StringValue("failed to convert token id into string")})

			render.JSON(w, r, resp.Error("failed to convert token id into string"))

			return
		}

		refreshToken, err := refreshTokener.GetRefreshTokenById(tokenIdStr)
		if err != nil {
			log.Error("failed to get user id", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to get user id"))

			return
		}

		hashedRequestRefreshToken := auth.HashRefreshToken(req.RefreshToken, cfg.Authorization.Salt)
		if refreshToken.Revoked {
			log.Error("token revoked", slog.Attr{Key: "error", Value: slog.StringValue("token revoked")})

			render.JSON(w, r, resp.RevokedRefreshToken())

			return
		}
		if refreshToken.ExpiresAt.Before(time.Now()) {
			log.Error("token expired", slog.Attr{Key: "error", Value: slog.StringValue("token expired")})

			_ = refreshTokener.RevokeRefreshTokenById(refreshToken.Id)

			render.JSON(w, r, resp.RevokedRefreshToken())

			return
		}
		if refreshToken.TokenHash != hashedRequestRefreshToken {
			log.Error("invalid refresh token", slog.Attr{Key: "error", Value: slog.StringValue("invalid refresh token")})

			render.JSON(w, r, resp.Error("invalid refresh token"))

			return
		}

		accessExp := time.Now().Add(cfg.Authorization.AccessTTL)

		_, accessToken, err := tokenAuth.Encode(map[string]interface{}{
			"user_id": refreshToken.UserId,
			"exp":     accessExp,
		})
		if err != nil {
			log.Error("failed to encode access token", slog.Attr{Key: "error", Value: slog.StringValue(err.Error())})

			render.JSON(w, r, resp.Error("failed to encode access token"))

			return
		}

		log.Info("token refreshed", slog.String("user_id", refreshToken.UserId))

		render.JSON(w, r, Response{
			resp.OK(),
			accessToken,
		})
	}
}
