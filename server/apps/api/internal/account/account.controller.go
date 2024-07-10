package account

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/middlewares/authentication"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/moonrhythm/validator"
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

func (a *createAccountRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := a.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (a *createAccountRequest) Validate() error {
	v := validator.New()
	v.Must(a.Type != "", "type is required")
	v.Must(a.Name != "", "name is required")
	v.Must(a.Bank != "", "bank is required")

	return errors.WithStack(v.Error())
}

func (h *controller) CreateAccount(ctx *fiber.Ctx) error {
	var req createAccountRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
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
	Id   uuid.UUID `params:"id"`
	Type string    `json:"type"`
	Name string    `json:"name"`
	Bank string    `json:"bank"`
}

func (a *updateAccountRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := ctx.BodyParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := a.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (a *updateAccountRequest) Validate() error {
	v := validator.New()
	v.Must(a.Id != uuid.Nil, "id is required")

	return errors.WithStack(v.Error())
}

func (h *controller) UpdateAccount(ctx *fiber.Ctx) error {
	var req updateAccountRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	account, err := h.usecase.UpdateAccount(ctx.UserContext(), userID, req.Id, entity.AccountInput{
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
	Id     uuid.UUID       `params:"id"`
	Amount decimal.Decimal `json:"amount"`
}

func (a *depositRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := ctx.BodyParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := a.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (a *depositRequest) Validate() error {
	v := validator.New()
	v.Must(a.Id != uuid.Nil, "id is required")
	v.Must(a.Amount.IsPositive(), "amount must be positive")

	return errors.WithStack(v.Error())
}

func (h *controller) Deposit(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	var req depositRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&dto.HttpResponse{
			Error: err.Error(),
		})
	}

	if err := h.usecase.Deposit(ctx.UserContext(), userID, req.Id, req.Amount); err != nil {
		return errors.Wrap(err, "failed to deposit")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "success",
	})
}
