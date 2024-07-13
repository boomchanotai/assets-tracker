package account

import (
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/middlewares/authentication"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/pocket"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/moonrhythm/validator"
	"github.com/shopspring/decimal"
)

type controller struct {
	usecase        *usecase
	pocketUsecase  *pocket.Usecase
	authMiddleware authentication.AuthMiddleware
}

func NewController(
	accountUsecase *usecase,
	pocketUsecase *pocket.Usecase,
	authMiddleware authentication.AuthMiddleware,
) *controller {
	return &controller{
		usecase:        accountUsecase,
		pocketUsecase:  pocketUsecase,
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
	ID        uuid.UUID               `json:"id"`
	UserID    uuid.UUID               `json:"userId"`
	Type      entity.AccountType      `json:"type"`
	Name      string                  `json:"name"`
	Bank      string                  `json:"bank"`
	Balance   decimal.Decimal         `json:"balance"`
	CreatedAt int64                   `json:"createdAt"`
	UpdatedAt int64                   `json:"updatedAt"`
	Pockets   []pocket.PocketResponse `json:"pockets"`
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
			CreatedAt: account.CreatedAt.Unix(),
			UpdatedAt: account.UpdatedAt.Unix(),
		})
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountsResponse,
	})
}

type accountRequest struct {
	Id uuid.UUID `params:"id"`
}

func (a *accountRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(a); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := a.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (a *accountRequest) Validate() error {
	v := validator.New()
	v.Must(a.Id != uuid.Nil, "id is required")

	return errors.WithStack(v.Error())
}

func (h *controller) GetAccount(ctx *fiber.Ctx) error {
	var req accountRequest
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

	account, err := h.usecase.GetAccount(ctx.UserContext(), userID, req.Id)
	if err != nil {
		return errors.Wrap(err, "failed to get account")
	}

	pockets, err := h.pocketUsecase.GetPocketsByAccountID(ctx.UserContext(), userID, req.Id)
	if err != nil {
		return errors.Wrap(err, "failed to get pockets")
	}

	pocketsResponse := make([]pocket.PocketResponse, 0, len(pockets))
	for _, p := range pockets {
		pocketsResponse = append(pocketsResponse, pocket.PocketResponse{
			ID:        p.ID,
			AccountID: p.AccountID,
			Name:      p.Name,
			Balance:   p.Balance,
			CreatedAt: p.CreatedAt.Unix(),
			UpdatedAt: p.UpdatedAt.Unix(),
		})
	}

	return ctx.JSON(dto.HttpResponse{
		Result: accountResponse{
			ID:        account.ID,
			UserID:    account.UserID,
			Type:      account.Type,
			Name:      account.Name,
			Bank:      account.Bank,
			Balance:   account.Balance,
			CreatedAt: account.CreatedAt.Unix(),
			UpdatedAt: account.UpdatedAt.Unix(),
			Pockets:   pocketsResponse,
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
			CreatedAt: account.CreatedAt.Unix(),
			UpdatedAt: account.UpdatedAt.Unix(),
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
			CreatedAt: account.CreatedAt.Unix(),
			UpdatedAt: account.UpdatedAt.Unix(),
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
