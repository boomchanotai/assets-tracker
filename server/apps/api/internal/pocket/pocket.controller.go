package pocket

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

func NewController(pocketUsecase *usecase, authMiddleware authentication.AuthMiddleware) *controller {
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

type pocketResponse struct {
	ID        uuid.UUID       `json:"id"`
	AccountID uuid.UUID       `json:"accountId"`
	Name      string          `json:"name"`
	Balance   decimal.Decimal `json:"balance"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
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

	res := make([]pocketResponse, 0, len(pockets))
	for _, pocket := range pockets {
		res = append(res, pocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt,
			UpdatedAt: pocket.UpdatedAt,
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
		Result: pocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt,
			UpdatedAt: pocket.UpdatedAt,
		},
	})
}

type createPocketRequest struct {
	AccountID string `json:"accountId"`
	Name      string `json:"name"`
}

func (h *controller) CreatePocket(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	var req createPocketRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if req.AccountID == "" || req.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	accountID, err := uuid.Parse(req.AccountID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	pocket, err := h.usecase.CreatePocket(ctx.UserContext(), entity.PocketInput{
		UserID:    userID,
		AccountID: accountID,
		Name:      req.Name,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: pocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt,
			UpdatedAt: pocket.UpdatedAt,
		},
	})
}

type updatePocketRequest struct {
	Name string `json:"name"`
}

func (h *controller) UpdatePocket(ctx *fiber.Ctx) error {
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

	var req updatePocketRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	pocket, err := h.usecase.UpdatePocket(ctx.UserContext(), userID, pocketID, entity.PocketInput{
		Name: req.Name,
	})
	if err != nil {
		return errors.Wrap(err, "failed to update pocket")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: pocketResponse{
			ID:        pocket.ID,
			AccountID: pocket.AccountID,
			Name:      pocket.Name,
			Balance:   pocket.Balance,
			CreatedAt: pocket.CreatedAt,
			UpdatedAt: pocket.UpdatedAt,
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
	ToPocketID string `json:"toPocketId"`
	Amount     string `json:"amount"`
}

func (h *controller) Transfer(ctx *fiber.Ctx) error {
	panic("not implemented")
}

type withdrawRequest struct {
	Amount string `json:"amount"`
}

func (h *controller) Withdraw(ctx *fiber.Ctx) error {
	panic("not implemented")
}
