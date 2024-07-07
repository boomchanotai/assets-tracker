package interfaces

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
)

type AccountRepository interface {
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error)
	GetUserAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Account, error)
	CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, input entity.AccountInput) (*entity.Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}
