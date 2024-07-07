package pocket

import (
	"context"
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/model"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("INSUFFICIENT_BALANCE")
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.PocketRepository {
	db.AutoMigrate(&model.Pocket{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetPocketsByAccountID(ctx context.Context, accountID uuid.UUID) ([]entity.Pocket, error) {
	var pockets []*model.Pocket
	if err := r.db.Where("account_id = ?", accountID).Find(&pockets).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get pockets")
	}

	var result []entity.Pocket
	for _, p := range pockets {
		result = append(result, entity.Pocket{
			ID:        p.ID,
			AccountID: p.AccountID,
			Name:      p.Name,
			Type:      p.Type,
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
		Type:      pocket.Type,
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
		Type:      input.Type,
		Balance:   decimal.NewFromInt(0), // Initial balance is 0
	}

	if err := r.db.Create(&p).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create pocket")
	}

	return &entity.Pocket{
		ID:        p.ID,
		AccountID: p.AccountID,
		Name:      p.Name,
		Type:      p.Type,
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
		Type:      p.Type,
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

// Deposit can be done only Cashbox pocket
func (r *repository) Deposit(ctx context.Context, pocketID uuid.UUID, amount decimal.Decimal) error {
	var pocket model.Pocket
	if err := r.db.First(&pocket, pocketID).Error; err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	if pocket.Type != entity.PocketTypeCashBox {
		return errors.Wrap(errors.New("INVALID_POCKET"), "failed to deposit")
	}

	// TODO: Lock db transaction
	pocket.Balance = pocket.Balance.Add(amount)
	pocket.UpdatedAt = time.Now()

	if err := r.db.Save(&pocket).Error; err != nil {
		return errors.Wrap(err, "failed to deposit")
	}

	return nil
}

func (r *repository) Transfer(ctx context.Context, fromPocketID, toPocketID uuid.UUID, amount decimal.Decimal) error {
	var fromPocket model.Pocket
	if err := r.db.First(&fromPocket, fromPocketID).Error; err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	var toPocket model.Pocket
	if err := r.db.First(&toPocket, toPocketID).Error; err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	if fromPocket.Balance.LessThan(amount) {
		return errors.Wrap(ErrInsufficientBalance, "failed to transfer")
	}

	fromPocket.Balance = fromPocket.Balance.Sub(amount)
	fromPocket.UpdatedAt = time.Now()

	toPocket.Balance = toPocket.Balance.Add(amount)
	toPocket.UpdatedAt = time.Now()

	tx := r.db.Begin()
	if err := tx.Save(&fromPocket).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to transfer")
	}

	if err := tx.Save(&toPocket).Error; err != nil {
		tx.Rollback()
		return errors.Wrap(err, "failed to transfer")
	}

	tx.Commit()

	return nil
}

func (r *repository) Withdraw(ctx context.Context, pocketID uuid.UUID, amount decimal.Decimal) error {
	var pocket model.Pocket
	if err := r.db.First(&pocket, pocketID).Error; err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	if pocket.Balance.LessThan(amount) {
		return errors.Wrap(ErrInsufficientBalance, "failed to withdraw")
	}

	pocket.Balance = pocket.Balance.Sub(amount)
	pocket.UpdatedAt = time.Now()

	if err := r.db.Save(&pocket).Error; err != nil {
		return errors.Wrap(err, "failed to withdraw")
	}

	return nil
}
