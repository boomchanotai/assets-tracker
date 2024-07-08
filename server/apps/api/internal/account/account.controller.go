package account

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
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

func NewController(accountUsecase *usecase, authMiddleware authentication.AuthMiddleware) *controller {
	return &controller{
		usecase:        accountUsecase,
		authMiddleware: authMiddleware,
	}
}

func (h *controller) Mount(r fiber.Router) {
	r.Get("/", h.GetAccounts)
	r.Get("/:id", h.GetAccount)
	r.Post("/", h.CreateAccount)
	r.Put("/:id", h.UpdateAccount)
	r.Delete("/:id", h.DeleteAccount)

	r.Post("/:id/deposit", h.Deposit)
	// r.Post("/:id/update-balance", h.Update)
}

type accountResponse struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"userId"`
	Type      entity.AccountType `json:"type"`
	Name      string             `json:"name"`
	Bank      string             `json:"bank"`
	Balance   decimal.Decimal    `json:"balance"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}

func (h *controller) GetAccounts(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	accounts, err := h.usecase.GetAccounts(ctx.UserContext(), userID)
	if err != nil {
		return errors.Wrap(err, "failed to get accounts")
	}

	accountsResponse := make([]accountResponse, 0, len(accounts))
	for _, account := range accounts {
		accountsResponse = append(accountsResponse, accountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Type:      account.Type,
			Name:      account.Name,
			Bank:      account.Bank,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		})
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountsResponse,
	})
}

func (h *controller) GetAccount(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	paramId := ctx.Params("id")
	accountId, err := uuid.Parse(paramId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	account, err := h.usecase.GetAccount(ctx.UserContext(), userID, accountId)
	if err != nil {
		return errors.Wrap(err, "failed to get account")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Type:      account.Type,
			Name:      account.Name,
			Bank:      account.Bank,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
	})
}

type createAccountRequest struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Bank string `json:"bank"`
}

func (h *controller) CreateAccount(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	var req createAccountRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if req.Type == "" || req.Name == "" || req.Bank == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	account, err := h.usecase.CreateAccount(ctx.UserContext(), entity.AccountInput{
		UserID: userID,
		Type:   entity.AccountType(req.Type),
		Name:   req.Name,
		Bank:   req.Bank,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create account")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Type:      account.Type,
			Name:      account.Name,
			Bank:      account.Bank,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
	})
}

type updateAccountRequest struct {
	Type    string          `json:"type"`
	Name    string          `json:"name"`
	Bank    string          `json:"bank"`
	Balance decimal.Decimal `json:"balance"`
}

func (h *controller) UpdateAccount(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	var req updateAccountRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	paramId := ctx.Params("id")
	accountId, err := uuid.Parse(paramId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	account, err := h.usecase.UpdateAccount(ctx.UserContext(), userID, accountId, entity.AccountInput{
		Type: entity.AccountType(req.Type),
		Name: req.Name,
		Bank: req.Bank,
	})

	if err != nil {
		return errors.Wrap(err, "failed to update account")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Type:      account.Type,
			Name:      account.Name,
			Bank:      account.Bank,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt,
			UpdatedAt: account.UpdatedAt,
		},
	})
}

func (h *controller) DeleteAccount(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	paramId := ctx.Params("id")
	accountId, err := uuid.Parse(paramId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	err = h.usecase.DeleteAccount(ctx.UserContext(), userID, accountId)
	if err != nil {
		return errors.Wrap(err, "failed to delete account")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "Account deleted",
	})
}

type depositRequest struct {
	Amount decimal.Decimal `json:"amount"`
}

func (h *controller) Deposit(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	var req depositRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if req.Amount.IsNegative() {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	paramId := ctx.Params("id")
	accountId, err := uuid.Parse(paramId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if err := h.usecase.Deposit(ctx.UserContext(), userID, accountId, req.Amount); err != nil {
		return errors.Wrap(err, "failed to deposit")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "success",
	})
}

// type updateRequest struct {
// 	Balance decimal.Decimal `json:"balance"`
// }

// func (h *controller) Update(ctx *fiber.Ctx) error {
// 	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
// 	if err != nil {
// 		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
// 			Error: "Unauthorized",
// 		})
// 	}

// 	var req updateRequest
// 	if err := ctx.BodyParser(&req); err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
// 			Error: "Bad Request",
// 		})
// 	}

// 	if req.Balance <= 0 {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
// 			Error: "Bad Request",
// 		})
// 	}

// 	paramId := ctx.Params("id")
// 	accountId, err := uuid.Parse(paramId)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
// 			Error: "Bad Request",
// 		})
// 	}

// 	account, err := h.usecase.UpdateBalance(ctx.UserContext(), userID, accountId, decimal.NewFromFloat(req.Balance))
// 	if err != nil {
// 		return errors.Wrap(err, "failed to update balance")
// 	}

// 	return ctx.JSON(dto.HttpResponse{
// 		Result: accountResponse{
// 			ID:        account.ID,
// 			UserID:    account.UserID,
// 			Type:      account.Type,
// 			Name:      account.Name,
// 			Bank:      account.Bank,
// 			Balance:   account.Balance,
// 			CreatedAt: account.CreatedAt,
// 			UpdatedAt: account.UpdatedAt,
// 		},
// 	})
// }
