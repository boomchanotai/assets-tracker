package interfaces

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
)

type TransactionRepository interface {
	GetTransactionByAccountID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) ([]entity.Transaction, error)
	CreateTransaction(ctx context.Context, transaction entity.TransactionInput) (*entity.Transaction, error)
}
