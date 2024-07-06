package middlewares

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	jwt "github.com/boomchanotai/assets-tracker/server/apps/api/internal/utils"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type UserId struct{}

func Auth(ctx *fiber.Ctx) error {
	tokenByte := ctx.GetReqHeaders()["Authorization"]
	if len(tokenByte) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "TOKEN_NOT_FOUND",
		})
	}

	bearerToken := tokenByte[0][7:]
	if len(bearerToken) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Result: "TOKEN_NOT_FOUND",
		})
	}

	// TODO: Implement Get Cache token
	cachedToken := entity.CachedTokens{}

	claims, err := validateToken(bearerToken, cachedToken)
	if err != nil {
		return errors.Wrap(err, "failed to validate token")
	}

	userContext := withUserID(ctx.UserContext(), claims.ID)
	ctx.SetUserContext(userContext)

	return ctx.Next()
}

func validateToken(accessToken string, cachedToken entity.CachedTokens) (*entity.JWTentity, error) {
	// TODO: SECRET should be stored in config
	claims, err := jwt.ParseToken(accessToken, "SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token")
	}

	err = jwt.ValidateToken(&cachedToken, claims, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate refresh token")
	}

	return claims, nil
}

type userIDContext struct{}

func withUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContext{}, userID)
}

func GetUserIDFromContext(ctx context.Context) (*uint, error) {
	userID, ok := ctx.Value(userIDContext{}).(uint)

	if !ok {
		return nil, errors.New("failed to get user id from context")
	}

	return &userID, nil
}
