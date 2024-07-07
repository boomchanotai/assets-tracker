package account

import (
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	usecase *usecase
}

func NewController(accountUsecase *usecase) *controller {
	return &controller{
		usecase: accountUsecase,
	}
}

func (h *controller) Mount(r fiber.Router) {
	r.Get("/", h.GetAccounts)
	r.Get("/:id", h.GetAccount)
	r.Post("/", h.CreateAccount)
	r.Put("/:id", h.UpdateAccount)
	r.Delete("/:id", h.DeleteAccount)
}

func (h *controller) GetAccounts(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "GetAccounts",
	})
}

func (h *controller) GetAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "GetAccount",
	})
}

func (h *controller) CreateAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "CreateAccount",
	})
}

func (h *controller) UpdateAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "UpdateAccount",
	})
}

func (h *controller) DeleteAccount(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "DeleteAccount",
	})
}
