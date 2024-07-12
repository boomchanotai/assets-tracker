package authentication

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/jwt"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("INVALID_TOKEN")
)

type AuthMiddleware interface {
	Auth(ctx *fiber.Ctx) error
	GetUserIDFromContext(ctx context.Context) (uuid.UUID, error)
}

type authMiddleware struct {
	userRepo interfaces.UserRepository
	config   *jwt.Config
}

func NewAuthMiddleware(userRepo interfaces.UserRepository, config *jwt.Config) AuthMiddleware {
	return &authMiddleware{
		userRepo: userRepo,
		config:   config,
	}
}

func (r *authMiddleware) Auth(ctx *fiber.Ctx) error {
	tokenByte := ctx.GetReqHeaders()["Authorization"]
	if len(tokenByte) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.HttpResponse{
			Result: "Unauthorized",
		})
	}

	if len(tokenByte[0]) < 7 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.HttpResponse{
			Result: "Unauthorized",
		})
	}

	bearerToken := tokenByte[0][7:]
	if len(bearerToken) == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.HttpResponse{
			Result: "Unauthorized",
		})
	}

	claims, err := r.validateToken(ctx.UserContext(), bearerToken)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.HttpResponse{
			Result: "Unauthorized",
		})
	}

	userContext := r.withUserID(ctx.UserContext(), claims.ID)
	ctx.SetUserContext(userContext)

	return ctx.Next()
}

func (r *authMiddleware) validateToken(ctx context.Context, bearerToken string) (*jwt.JWTentity, error) {
	parsedToken, err := jwt.ParseToken(bearerToken, r.config.AccessTokenSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token")
	}

	cachedToken, err := r.userRepo.GetUserAuthToken(ctx, parsedToken.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get cached token")
	}

	err = jwt.ValidateToken(cachedToken, parsedToken, false)
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate refresh token")
	}

	return parsedToken, nil

}

type userIDContext struct{}

func (r *authMiddleware) withUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContext{}, userID)
}

func (r *authMiddleware) GetUserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(userIDContext{}).(uuid.UUID)

	if !ok {
		return uuid.UUID{}, errors.New("failed to get user id from context")
	}

	return userID, nil
}
