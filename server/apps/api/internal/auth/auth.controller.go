package auth

import (
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/middlewares"
	"github.com/cockroachdb/errors"
	"github.com/gofiber/fiber/v2"
)

func (h *controller) Mount(r fiber.Router, authMiddleware middlewares.AuthMiddleware) {
	r.Post("/register", h.Register)
	r.Post("/login", h.Login)
	r.Get("/me", authMiddleware.Auth, h.GetProfile)
}

type controller struct {
	usecase        *usecase
	authMiddleware middlewares.AuthMiddleware
}

func NewController(authUsecase *usecase, authMiddleware middlewares.AuthMiddleware) *controller {
	return &controller{
		usecase:        authUsecase,
		authMiddleware: authMiddleware,
	}
}

type registerRequest struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type registerResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

type getProfileReponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *controller) GetProfile(ctx *fiber.Ctx) error {
	userID, err := h.authMiddleware.GetUserIDFromContext(ctx.UserContext())
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(&dto.HttpResponse{
			Error: "Unauthorized",
		})
	}

	user, err := h.usecase.GetProfile(ctx.UserContext(), userID)
	if err != nil {
		return errors.Wrap(err, "failed to get user")
	}

	return ctx.JSON(dto.HttpResponse{
		Result: getProfileReponse{
			ID:    user.ID.String(),
			Email: user.Email,
			Name:  user.Name,
		},
	})
}
