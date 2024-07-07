package transaction

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/model"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetTransactionsByPocketID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) ([]entity.Transaction, error)
	GetTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*entity.Transaction, error)
	CreateTransaction(ctx context.Context, transaction entity.TransactionInput) (*entity.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&model.Transaction{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetTransactionsByPocketID(ctx context.Context, userID uuid.UUID, pocketID uuid.UUID) ([]entity.Transaction, error) {
	var transactions []*model.Transaction
	if err := r.db.Where("user_id = ? AND (from_pocket_id = ? OR to_pocket_id = ?)", userID, pocketID, pocketID).Find(&transactions).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get transactions")
	}

	var result []entity.Transaction
	for _, t := range transactions {
		result = append(result, entity.Transaction{
			ID:           t.ID,
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

func (r *repository) GetTransaction(ctx context.Context, userID uuid.UUID, transactionID uuid.UUID) (*entity.Transaction, error) {
	var t model.Transaction
	if err := r.db.Where("user_id = ? AND id = ?", userID, transactionID).First(&t).Error; err != nil {
		return nil, errors.Wrap(err, "failed to get transaction")
	}

	return &entity.Transaction{
		ID:           t.ID,
		FromPocketID: t.FromPocketID,
		ToPocketID:   t.ToPocketID,
		Type:         entity.TxType(t.Type),
		Amount:       t.Amount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}, nil
}

func (r *repository) CreateTransaction(ctx context.Context, transaction entity.TransactionInput) (*entity.Transaction, error) {

	t := model.Transaction{
		UserID:       transaction.UserID,
		FromPocketID: transaction.FromPocketID,
		ToPocketID:   transaction.ToPocketID,
		Type:         transaction.Type,
		Amount:       transaction.Amount,
	}

	if err := r.db.Create(&t).Error; err != nil {
		return nil, errors.Wrap(err, "failed to create transaction")
	}

	return &entity.Transaction{
		ID:           t.ID,
		FromPocketID: t.FromPocketID,
		ToPocketID:   t.ToPocketID,
		Type:         entity.TxType(t.Type),
		Amount:       t.Amount,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
	}, nil
}
