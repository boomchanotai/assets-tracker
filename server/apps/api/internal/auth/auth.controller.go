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

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

func (h *controller) Register(ctx *fiber.Ctx) error {
	var body registerRequest
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

	res, err := h.usecase.Register(ctx.Context(), body.Email, body.Name, body.Password)
	if errors.Is(err, ErrEmailAlreadyExists) {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Email already exists",
		})
	}
	if err != nil {
		return errors.Wrap(err, "failed to register")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: registerResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Exp:          res.Exp,
		},
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Exp          int64  `json:"exp"`
}

func (h *controller) Login(ctx *fiber.Ctx) error {
	var body loginRequest
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

	if body.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(dto.HttpResponse{
			Error: "Password is required",
		})
	}

	res, err := h.usecase.Login(ctx.UserContext(), body.Email, body.Password)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	return ctx.JSON(dto.HttpResponse{
		Result: loginResponse{
			AccessToken:  res.AccessToken,
			RefreshToken: res.RefreshToken,
			Exp:          res.Exp,
		},
	})
}
