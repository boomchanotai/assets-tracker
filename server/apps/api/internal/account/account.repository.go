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
	GetAccounts(ctx context.Context) ([]*entity.Account, error)
	GetAccount(ctx context.Context, id string) (*entity.Account, error)
	CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error)
	UpdateAccount(ctx context.Context, id string, input entity.AccountInput) (*entity.Account, error)
	DeleteAccount(ctx context.Context, id string) error
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

func (r *repository) GetAccounts(ctx context.Context) ([]*entity.Account, error) {
	var accounts []*model.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get accounts")
	}

	var result []*entity.Account
	for _, a := range accounts {
		result = append(result, &entity.Account{
			ID:        a.ID,
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

func (r *repository) GetAccount(ctx context.Context, id string) (*entity.Account, error) {
	var a model.Account
	if err := r.db.Where("id = ?", id).First(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get account")
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

func (r *repository) CreateAccount(ctx context.Context, input entity.AccountInput) (*entity.Account, error) {
	a := model.Account{
		ID:        uuid.New(),
		Type:      input.Type,
		Name:      input.Name,
		Bank:      input.Bank,
		Balance:   decimal.NewFromInt(0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create account")
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

func (r *repository) UpdateAccount(ctx context.Context, id string, input entity.AccountInput) (*entity.Account, error) {
	var a model.Account
	if err := r.db.Where("id = ?", id).First(&a).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get account")
	}

	a.Type = input.Type
	a.Name = input.Name
	a.Bank = input.Bank
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

func (r *repository) DeleteAccount(ctx context.Context, id string) error {
	if err := r.db.Where("id = ?", id).Delete(&model.Account{}).Error; err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return nil
}
