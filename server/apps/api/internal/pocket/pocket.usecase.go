package pocket

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type usecase struct {
	pocketRepo  interfaces.PocketRepository
	accountRepo interfaces.AccountRepository
}

func NewUsecase(pocketRepo interfaces.PocketRepository, accountRepo interfaces.AccountRepository) *usecase {
	return &usecase{
		pocketRepo:  pocketRepo,
		accountRepo: accountRepo,
	}
}

func (u *usecase) GetPockets(ctx context.Context, userID uuid.UUID) ([]entity.Pocket, error) {
	pockets, err := u.pocketRepo.GetPocketByUserID(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pockets")
	}

	return pockets, nil
}

func (u *usecase) GetPocket(ctx context.Context, userID, pocketID uuid.UUID) (*entity.Pocket, error) {
	pocket, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pocket")
	}

	return pocket, nil
}

func (u *usecase) CreatePocket(ctx context.Context, input entity.PocketInput) (*entity.Pocket, error) {
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

func (u *usecase) UpdatePocket(ctx context.Context, userID, pocketID uuid.UUID, input entity.PocketInput) (*entity.Pocket, error) {
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

func (u *usecase) DeletePocket(ctx context.Context, userID, pocketID uuid.UUID) error {
	// Check ownership
	if _, err := u.pocketRepo.GetPocketByID(ctx, userID, pocketID); err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	if err := u.pocketRepo.DeletePocket(ctx, pocketID); err != nil {
		return errors.Wrap(err, "failed to delete pocket")
	}

	return nil
}
