package transaction

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/model"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) interfaces.TransactionRepository {
	db.AutoMigrate(&model.Transaction{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetTransactionByAccountID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) ([]entity.Transaction, error) {
	var transactions []*model.Transaction
	if err := r.db.Where("account_id IN (?)", r.db.Model(&model.Account{}).Select("id").Where("user_id = ?", userID)).Order("created_at desc").Find(&transactions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get transactions")
	}

	var result []entity.Transaction
	for _, t := range transactions {
		result = append(result, entity.Transaction{
			ID:           t.ID,
			AccountID:    t.AccountID,
			FromPocketID: t.FromPocketID,
			ToPocketID:   t.ToPocketID,
			Type:         entity.TxType(t.Type),
			Amount:       t.Amount,
			CreatedAt:    t.CreatedAt,
			UpdatedAt:    t.UpdatedAt,
		})
	}

	return result, nil
}

func (r *repository) CreateTransaction(ctx context.Context, input entity.TransactionInput) (*entity.Transaction, error) {
	t := model.Transaction{
		ID:           uuid.New(),
		AccountID:    input.AccountID,
		FromPocketID: input.FromPocketID,
		ToPocketID:   input.ToPocketID,
		Type:         input.Type,
		Amount:       input.Amount,
	}

	if err := r.db.Create(&t).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}

	return &entity.Transaction{
		ID:           t.ID,
		AccountID:    t.AccountID,
		FromPocketID: t.FromPocketID,
		ToPocketID:   t.ToPocketID,
		Type:         entity.TxType(t.Type),
		Amount:       t.Amount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}, nil
}
