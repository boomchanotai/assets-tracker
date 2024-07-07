package transaction

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/middlewares/authentication"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type controller struct {
	usecase        *usecase
	authMiddleware authentication.AuthMiddleware
}

func NewController(usecase *usecase, authMiddleware authentication.AuthMiddleware) *controller {
	return &controller{
		usecase:        usecase,
		authMiddleware: authMiddleware,
	}
}

func (h *controller) Mount(r fiber.Router) {
	r.Get("/:id", h.GetTransactionByAccountID)
}

type transactionResponse struct {
	ID           uuid.UUID       `json:"id"`
	AccountID    uuid.UUID       `json:"accountId"`
	FromPocketID *uuid.UUID      `json:"fromPocketId"`
	ToPocketID   *uuid.UUID      `json:"toPocketId"`
	Amount       decimal.Decimal `json:"amount"`
	CreatedAt    time.Time       `json:"createdAt"`
	UpdatedAt    time.Time       `json:"updatedAt"`
}

func (h *controller) GetTransactionByAccountID(ctx *fiber.Ctx) error {
	// Get user ID from context
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	// Get account ID from path parameter
	accountID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Invalid account ID",
		})
	}

	// Get transactions
	transactions, err := h.usecase.GetTransactions(ctx.UserContext(), userID, accountID)
	if err != nil {
		return errors.Wrap(err, "failed to get transactions")
	}

	// Response
	res := make([]transactionResponse, 0, len(transactions))
	for _, transaction := range transactions {
		res = append(res, transactionResponse{
			ID:           transaction.ID,
			AccountID:    transaction.AccountID,
			FromPocketID: transaction.FromPocketID,
			ToPocketID:   transaction.ToPocketID,
			Amount:       transaction.Amount,
			CreatedAt:    transaction.CreatedAt,
			UpdatedAt:    transaction.UpdatedAt,
		})
	}

	return ctx.JSON(&dto.HttpResponse{
		Result: res,
	})
}
