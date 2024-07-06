package user

import (
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/gofiber/fiber/v2"
)

func (h *controller) Mount(r fiber.Router) {
	r.Get("/", h.GetUsers)
	r.Get("/:id", h.GetUser)
}

type controller struct {
	usecase *usecase
}

func NewController(userUsecase *usecase) *controller {
	return &controller{
		usecase: userUsecase,
	}
}

func (h *controller) GetUsers(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "GetUsers",
	})
}

func (h *controller) GetUser(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "GetUser",
	})
}
