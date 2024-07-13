package pocket

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type Usecase struct {
	pocketRepo      interfaces.PocketRepository
	accountRepo     interfaces.AccountRepository
	transactionRepo interfaces.TransactionRepository
}

func NewUsecase(pocketRepo interfaces.PocketRepository, accountRepo interfaces.AccountRepository, transactionRepo interfaces.TransactionRepository) *Usecase {
	return &Usecase{
		pocketRepo:      pocketRepo,
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func (u *Usecase) GetPocket(ctx context.Context, userID, pocketID uuid.UUID) (*entity.Pocket, error) {
	pocket, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pocket")
	}

	return pocket, nil
}

func (u *Usecase) GetPocketsByAccountID(ctx context.Context, userID uuid.UUID, accountID uuid.UUID) ([]entity.Pocket, error) {
	// Check account ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, userID, accountID); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	pockets, err := u.pocketRepo.GetPocketsByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pockets")
	}

	return pockets, nil
}

func (u *Usecase) CreatePocket(ctx context.Context, input entity.PocketInput) (*entity.Pocket, error) {
	// Check account ownership
	if _, err := u.accountRepo.GetUserAccount(ctx, input.UserID, input.AccountID); err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	pocket, err := u.pocketRepo.CreatePocket(ctx, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create pocket")
	}

	return pocket, nil
}

func (u *Usecase) UpdatePocket(ctx context.Context, userID, pocketID uuid.UUID, input entity.PocketInput) (*entity.Pocket, error) {
	// Check ownership
	if _, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID); err != nil {
		return nil, errors.Wrap(err, "failed to get pocket")
	}

	pocket, err := u.pocketRepo.UpdatePocket(ctx, pocketID, input)
	if err != nil {
		return nil, errors.Wrap(err, "failed to update pocket")
	}

	return pocket, nil
}

func (u *Usecase) DeletePocket(ctx context.Context, userID, pocketID uuid.UUID) error {
	// Check ownership
	if _, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID); err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	if err := u.pocketRepo.DeletePocket(ctx, pocketID); err != nil {
		return errors.Wrap(err, "failed to delete pocket")
	}

	return nil
}

func (u *Usecase) Transfer(ctx context.Context, userID, fromPocketID, toPocketID uuid.UUID, amount decimal.Decimal) error {
	// Check ownership
	fromPocket, err := u.pocketRepo.GetPocketByID(ctx, userID, fromPocketID)
	if err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	toPocket, err := u.pocketRepo.GetPocketByID(ctx, userID, toPocketID)
	if err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	// Create transaction
	if _, err := u.transactionRepo.CreateTransaction(ctx, entity.TransactionInput{
		AccountID:    fromPocket.AccountID,
		FromPocketID: &fromPocket.ID,
		ToPocketID:   &toPocket.ID,
		Type:         entity.TxTypeDeposit,
		Amount:       amount,
	}); err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}

	if err := u.pocketRepo.Transfer(ctx, fromPocketID, toPocketID, amount); err != nil {
		return errors.Wrap(err, "failed to transfer")
	}

	return nil
}

func (u *Usecase) Withdraw(ctx context.Context, userID, pocketID uuid.UUID, amount decimal.Decimal) error {
	// Check ownership
	fromPocket, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID)
	if err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	// Create transaction
	if _, err := u.transactionRepo.CreateTransaction(ctx, entity.TransactionInput{
		AccountID:    fromPocket.AccountID,
		FromPocketID: &fromPocket.ID,
		ToPocketID:   nil,
		Type:         entity.TxTypeDeposit,
		Amount:       amount,
	}); err != nil {
		return errors.Wrap(err, "failed to create transaction")
	}

	if err := u.pocketRepo.Withdraw(ctx, pocketID, amount); err != nil {
		return errors.Wrap(err, "failed to withdraw")
	}

	return nil
}
