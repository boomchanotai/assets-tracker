package transaction

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type usecase struct {
	transactionRepo interfaces.TransactionRepository
	accountRepo     interfaces.AccountRepository
}

func NewUsecase(transactionRepo interfaces.TransactionRepository, accountRepo interfaces.AccountRepository) *usecase {
	return &usecase{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (u *usecase) GetTransactions(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]entity.Transaction, error) {
	// Check account ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, accountID); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	transactions, err := u.transactionRepo.GetTransactionByAccountID(ctx, userID, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get transactions")
	}

	return transactions, nil
}
