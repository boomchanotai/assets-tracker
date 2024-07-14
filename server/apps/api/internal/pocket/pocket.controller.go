package pocket

import (
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
	usecase        *Usecase
	authMiddleware authentication.AuthMiddleware
}

func NewController(pocketUsecase *Usecase, authMiddleware authentication.AuthMiddleware) *controller {
	return &controller{
		usecase:        pocketUsecase,
		authMiddleware: authMiddleware,
	}
}

func (h *controller) Mount(r fiber.Router) {
	r.Get("/account/:id", h.GetPocketsByAccountID)
	r.Get("/:id", h.GetPocket)
	r.Post("/", h.CreatePocket)
	r.Put("/:id", h.UpdatePocket)
	r.Delete("/:id", h.DeletePocket)

	r.Post("/:id/transfer", h.Transfer)
	r.Post("/:id/withdraw", h.Withdraw)
}

type PocketResponse struct {
	ID        uuid.UUID         `json:"id"`
	AccountID uuid.UUID         `json:"accountId"`
	Name      string            `json:"name"`
	Type      entity.PocketType `json:"type"`
	Balance   decimal.Decimal   `json:"balance"`
	CreatedAt int64             `json:"createdAt"`
	UpdatedAt int64             `json:"updatedAt"`
}

func (h *controller) GetPocketsByAccountID(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	accountID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	pockets, err := h.usecase.GetPocketsByAccountID(ctx.UserContext(), userID, accountID)
	if err != nil {
		return errors.Wrap(err, "failed to get pockets")
	}

	res := make([]PocketResponse, 0, len(pockets))
	for _, pocket := range pockets {
		res = append(res, PocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Type:      pocket.Type,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt.Unix(),
			UpdatedAt: pocket.UpdatedAt.Unix(),
		})
	}

	return ctx.JSON(dto.HttpResponse{
		Result: res,
	})
}

func (h *controller) GetPocket(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	pocketID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	pocket, err := h.usecase.GetPocket(ctx.UserContext(), userID, pocketID)
	if err != nil {
		return errors.Wrap(err, "failed to get pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: PocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Type:      pocket.Type,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt.Unix(),
			UpdatedAt: pocket.UpdatedAt.Unix(),
		},
	})
}

type createPocketRequest struct {
	AccountID uuid.UUID `json:"accountId"`
	Name      string    `json:"name"`
}

func (p *createPocketRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.BodyParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := p.Validate(); err != nil {
		return errors.Wrap(err, "failed to validate request")
	}

	return nil
}

func (p *createPocketRequest) Validate() error {
	v := validator.New()
	v.Must(p.AccountID != uuid.Nil, "accountId is required")
	v.Must(p.Name != "", "name is required")

	return errors.WithStack(v.Error())
}

func (h *controller) CreatePocket(ctx *fiber.Ctx) error {
	var req createPocketRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	pocket, err := h.usecase.CreatePocket(ctx.UserContext(), entity.PocketInput{
		UserID:    userID,
		AccountID: req.AccountID,
		Name:      req.Name,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: PocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Type:      pocket.Type,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt.Unix(),
			UpdatedAt: pocket.UpdatedAt.Unix(),
		},
	})
}

type updatePocketRequest struct {
	Id   uuid.UUID `params:"id"`
	Name string    `json:"name"`
}

func (p *updatePocketRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := ctx.BodyParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := p.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (p *updatePocketRequest) Validate() error {
	v := validator.New()
	v.Must(p.Id != uuid.Nil, "id is required")
	v.Must(p.Name != "", "name is required")

	return errors.WithStack(v.Error())
}

func (h *controller) UpdatePocket(ctx *fiber.Ctx) error {
	var req updatePocketRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	pocket, err := h.usecase.UpdatePocket(ctx.UserContext(), userID, req.Id, entity.PocketInput{
		Name: req.Name,
	})
	if err != nil {
		return errors.Wrap(err, "failed to update pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: PocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Type:      pocket.Type,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt.Unix(),
			UpdatedAt: pocket.UpdatedAt.Unix(),
		},
	})
}

func (h *controller) DeletePocket(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	pocketID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if err := h.usecase.DeletePocket(ctx.UserContext(), userID, pocketID); err != nil {
		return errors.Wrap(err, "failed to delete pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "success",
	})
}

type transferRequest struct {
	FromPocketID uuid.UUID       `params:"id"`
	ToPocketID   uuid.UUID       `json:"toPocketId"`
	Amount       decimal.Decimal `json:"amount"`
}

func (p *transferRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := ctx.BodyParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := p.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (p *transferRequest) Validate() error {
	v := validator.New()
	v.Must(p.FromPocketID != uuid.Nil, "fromPocketId is required")
	v.Must(p.ToPocketID != uuid.Nil, "toPocketId is required")
	v.Must(!p.Amount.IsZero(), "amount is required")

	return errors.WithStack(v.Error())
}

func (h *controller) Transfer(ctx *fiber.Ctx) error {
	var req transferRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	if err := h.usecase.Transfer(ctx.UserContext(), userID, req.FromPocketID, req.ToPocketID, req.Amount); err != nil {
		return errors.Wrap(err, "failed to transfer")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "success",
	})
}

type withdrawRequest struct {
	Id     uuid.UUID       `params:"id"`
	Amount decimal.Decimal `json:"amount"`
}

func (p *withdrawRequest) Parse(ctx *fiber.Ctx) error {
	if err := ctx.ParamsParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := ctx.BodyParser(p); err != nil {
		return errors.Wrap(err, "failed to parse request")
	}

	if err := p.Validate(); err != nil {
		return errors.Wrap(err, "invalid request")
	}

	return nil
}

func (p *withdrawRequest) Validate() error {
	v := validator.New()
	v.Must(p.Id != uuid.Nil, "id is required")
	v.Must(!p.Amount.IsZero(), "amount is required")

	return errors.WithStack(v.Error())
}

func (h *controller) Withdraw(ctx *fiber.Ctx) error {
	var req withdrawRequest
	if err := req.Parse(ctx); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: err.Error(),
		})
	}

	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	if err := h.usecase.Withdraw(ctx.UserContext(), userID, req.Id, req.Amount); err != nil {
		return errors.Wrap(err, "failed to withdraw")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: "success",
	})
}
