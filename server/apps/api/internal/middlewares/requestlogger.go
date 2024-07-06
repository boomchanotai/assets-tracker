package middlewares

import (
	"log/slog"

	"github.com/boomchanotai/assets-tracker/server/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestLogger() func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.GetRespHeader("X-Request-Id")
		logger.InfoContext(ctx.UserContext(), "request received", slog.String("request_id", requestId), slog.String("method", ctx.Method()), slog.String("path", ctx.Path()))
		return ctx.Next()
	}
}
