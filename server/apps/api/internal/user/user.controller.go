package user

import (
	"github.com/boomchanotai/assets-tracker/server/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func (h *controller) Mount(r fiber.Router) {
	defer logger.Info("Mounted `block` api handlers")

	r.Get("/", h.GetUsers)
	r.Get("/:id", h.GetUser)
}

type controller struct {
	userUsecase *usecase
}

func NewController(userUsecase *usecase) *controller {
	return &controller{
		userUsecase: userUsecase,
	}
}

func (h *controller) GetUsers(c *fiber.Ctx) error {
	return c.SendString("GetUsers")
}

func (h *controller) GetUser(c *fiber.Ctx) error {
	return c.SendString("GetUser")
}
