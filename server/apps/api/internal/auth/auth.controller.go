package auth

import (
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

func (h *controller) Mount(r fiber.Router) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
}

type controller struct {
	usecase *usecase
}

func NewController(authUsecase *usecase) *controller {
	return &controller{
		usecase: authUsecase,
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

func (h *controller) Register(ctx *fiber.Ctx) error {
	var body RegisterRequest
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Bad Request",
		})
	}

	if body.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Email is required",
		})
	}

	if body.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Name is required",
		})
	}

	if body.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Password is required",
		})
	}

	user, err := h.usecase.Register(ctx.Context(), body.Email, body.Name, body.Password)
	if err != nil {
		return errors.Wrap(err, "failed to register")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: user,
	})
}

func (h *controller) Login(ctx *fiber.Ctx) error {
	return ctx.JSON(dto.HttpResponse{
		Result: "Login",
	})
}
