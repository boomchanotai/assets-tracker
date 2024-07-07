package interfaces

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, input entity.UserInput) (*entity.User, error)
	UpdateUser(ctx context.Context, id uuid.UUID, input entity.UserInput) (*entity.User, error)

	SetUserAuthToken(ctx context.Context, userID uuid.UUID, token entity.CachedTokens) error
	GetUserAuthToken(ctx context.Context, userID uuid.UUID) (*entity.CachedTokens, error)
	DeleteUserAuthToken(ctx context.Context, userID uuid.UUID) error
}
