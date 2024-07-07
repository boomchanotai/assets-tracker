package account

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type usecase struct {
	accountRepo Repository
}

func NewUsecase(accountRepo Repository) *usecase {
	return &usecase{
		accountRepo: accountRepo,
	}
}

func (u *usecase) GetAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error) {
	accounts, err := u.accountRepo.GetUserAccounts(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	return accounts, nil
}

func (u *usecase) GetAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Account, error) {
	account, err := u.accountRepo.GetUserAccount(ctx, userID, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	return account, nil
}

func (u *usecase) CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error) {
	account, err := u.accountRepo.CreateAccount(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	return account, nil
}

func (u *usecase) UpdateAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID, input entity.AccountInput) (*entity.Account, error) {
	account, err := u.accountRepo.UpdateAccount(ctx, userID, id, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update account")
	}

	return account, nil
}

func (u *usecase) DeleteAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	err := u.accountRepo.DeleteAccount(ctx, userID, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return nil
}
