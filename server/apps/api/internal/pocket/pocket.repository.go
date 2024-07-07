package pocket

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
	GetPocketByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Pocket, error)
	GetPocketByID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) (*entity.Pocket, error)
	CreatePocket(ctx context.Context, input entity.PocketInput) (*entity.Pocket, error)
	UpdatePocket(ctx context.Context, id uuid.UUID, input entity.PocketInput) (*entity.Pocket, error)
	DeletePocket(ctx context.Context, pocketID uuid.UUID) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&model.Pocket{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetPocketByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Pocket, error) {
	var pockets []*model.Pocket
	// where userID == pocket.Account.UserID
	if err := r.db.Where("account_id IN (?)", r.db.Model(&model.Account{}).Select("id").Where("user_id = ?", userID)).Find(&pockets).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get pockets")
	}

	var result []entity.Pocket
	for _, p := range pockets {
		result = append(result, entity.Pocket{
			ID:        p.ID,
			AccountID: p.AccountID,
			Name:      p.Name,
			Balance:   p.Balance,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	return result, nil
}

func (r *repository) GetPocketByID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) (*entity.Pocket, error) {
	var pocket model.Pocket
	// where userID == pocket.Account.UserID
	if err := r.db.Where("account_id IN (?)", r.db.Model(&model.Account{}).Select("id").Where("user_id = ?", userID)).First(&pocket, pocketID).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get pocket")
	}

	return &entity.Pocket{
		ID:        pocket.ID,
		AccountID: pocket.AccountID,
		Name:      pocket.Name,
		Balance:   pocket.Balance,
		CreatedAt: pocket.CreatedAt,
		UpdatedAt: pocket.UpdatedAt,
	}, nil
}

func (r *repository) CreatePocket(ctx context.Context, input entity.PocketInput) (*entity.Pocket, error) {
	p := model.Pocket{
		ID:        uuid.New(),
		AccountID: input.AccountID,
		Name:      input.Name,
		Balance:   decimal.NewFromInt(0), // Initial balance is 0
	}

	if err := r.db.Create(&p).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create pocket")
	}

	return &entity.Pocket{
		ID:        p.ID,
		AccountID: p.AccountID,
		Name:      p.Name,
		Balance:   p.Balance,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func (r *repository) UpdatePocket(ctx context.Context, id uuid.UUID, input entity.PocketInput) (*entity.Pocket, error) {
	var p model.Pocket
	if err := r.db.First(&p, id).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get pocket")
	}

	if input.Name != "" {
		p.Name = input.Name
	}

	p.UpdatedAt = time.Now()

	if err := r.db.Save(&p).Error; err != nil {
		return nil, err
	}

	return &entity.Pocket{
		ID:        p.ID,
		AccountID: p.AccountID,
		Name:      p.Name,
		Balance:   p.Balance,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}, nil
}

func (r *repository) DeletePocket(ctx context.Context, pocketID uuid.UUID) error {
	if err := r.db.Delete(&model.Pocket{}, pocketID).Error; err != nil {
		return errors.Wrap(err, "failed to delete pocket")
	}

	return nil
}
