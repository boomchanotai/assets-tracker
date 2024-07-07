package account

import (
	"context"
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/model"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Repository interface {
	GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error)
	GetUserAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Account, error)
	CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error)
	UpdateAccount(ctx context.Context, id uuid.UUID, input entity.AccountInput) (*entity.Account, error)
	DeleteAccount(ctx context.Context, id uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&model.Account{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetUserAccounts(ctx context.Context, userID uuid.UUID) ([]entity.Account, error) {
	var accounts []*model.Account
	if err := r.db.Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user accounts")
	}

	var result []entity.Account
	for _, a := range accounts {
		result = append(result, entity.Account{
			ID:        a.ID,
			UserID:    a.UserID,
			Type:      a.Type,
			Name:      a.Name,
			Bank:      a.Bank,
			Balance:   a.Balance,
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		})
	}

	return result, nil
}

func (r *repository) GetUserAccount(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*entity.Account, error) {
	var a model.Account
	if err := r.db.Where("user_id = ? AND id = ?", userID, id).First(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get user account")
	}

	return &entity.Account{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      a.Type,
		Name:      a.Name,
		Bank:      a.Bank,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}, nil
}

func (r *repository) CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error) {
	a := model.Account{
		ID:      uuid.New(),
		UserID:  input.UserID,
		Type:    input.Type,
		Name:    input.Name,
		Bank:    input.Bank,
		Balance: decimal.NewFromInt(0),
	}

	if err := r.db.Create(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create account")
	}

	return &entity.Account{
		ID:        a.ID,
		UserID:    a.UserID,
		Type:      a.Type,
		Name:      a.Name,
		Bank:      a.Bank,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}, nil
}

func (r *repository) UpdateAccount(ctx context.Context, id uuid.UUID, input entity.AccountInput) (*entity.Account, error) {
	var a model.Account
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	if input.Type != "" {
		a.Type = input.Type
	}

	if input.Name != "" {
		a.Name = input.Name
	}

	if input.Bank != "" {
		a.Bank = input.Bank
	}

	if input.Balance.IsPositive() {
		a.Balance = input.Balance
	}

	a.UpdatedAt = time.Now()

	if err := r.db.Save(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to update account")
	}

	return &entity.Account{
		ID:        a.ID,
		Type:      a.Type,
		Name:      a.Name,
		Bank:      a.Bank,
		Balance:   a.Balance,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}, nil
}

func (r *repository) DeleteAccount(ctx context.Context, id uuid.UUID) error {
	if err := r.db.Where("id = ?", id).Delete(&model.Account{}).Error; err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return nil
}
