package interfaces

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
)

type PocketRepository interface {
	GetPocketByID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) (*entity.Pocket, error)
	GetPocketsByAccountID(ctx context.Context, accountID uuid.UUID) ([]entity.Pocket, error)
	CreatePocket(ctx context.Context, input entity.PocketInput) (*entity.Pocket, error)
	UpdatePocket(ctx context.Context, id uuid.UUID, input entity.PocketInput) (*entity.Pocket, error)
	DeletePocket(ctx context.Context, pocketID uuid.UUID) error
}
