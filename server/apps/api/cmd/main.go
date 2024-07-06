package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/auth"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/config"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/dto"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/middlewares"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/user"
	"github.com/boomchanotai/assets-tracker/server/pkg/logger"
	"github.com/boomchanotai/assets-tracker/server/pkg/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	conf := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := logger.Init(conf.Logger); err != nil {
		logger.PanicContext(ctx, "failed to initialize logger", slog.Any("error", err))
	}

	db, err := gorm.Open(postgres.Open(conf.Postgres.String()), &gorm.Config{})
	if err != nil {
		logger.PanicContext(ctx, "failed to connect to database", slog.Any("error", err))
	}

	redisConn, err := redis.New(ctx, conf.Redis)
	if err != nil {
		logger.PanicContext(ctx, "failed to connect to redis", slog.Any("error", err))
	}
	defer func() {
		if err := redisConn.Close(); err != nil {
			logger.ErrorContext(ctx, "failed to close redis connection", slog.Any("error", err))
		}
	}()

	userRepo := user.NewRepository(db, redisConn, &conf.JWT)
	userUsecase := user.NewUsecase(userRepo)
	userController := user.NewController(userUsecase)

	authMiddleware := middlewares.NewAuthMiddleware(userRepo, &conf.JWT)

	authUsecase := auth.NewUsecase(userRepo, &conf.JWT)
	authController := auth.NewController(authUsecase, authMiddleware)

	app := fiber.New(fiber.Config{
		AppName:       conf.Name,
		CaseSensitive: true,
		JSONEncoder:   json.Marshal,
		JSONDecoder:   json.Unmarshal,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			logger.ErrorContext(c.UserContext(), "unhandled error", slog.Any("error", err))
			return c.Status(fiber.StatusInternalServerError).JSON(dto.HttpResponse{
				Error: "Internal Server Error",
			})
		},
	})

	app.Use(cors.New())

	authGroup := app.Group("/v1/auth")
	authController.Mount(authGroup, authMiddleware)

	userGroup := app.Group("/v1/user")
	userController.Mount(userGroup)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%d", conf.Port)); err != nil {
			logger.PanicContext(ctx, "failed to start server", slog.Any("error", err))
			stop()
		}
	}()

	defer func() {
		if err := app.ShutdownWithContext(ctx); err != nil {
			logger.ErrorContext(ctx, "failed to shutdown server", slog.Any("error", err))
		}
		logger.InfoContext(ctx, "gracefully shutdown server")
	}()

	<-ctx.Done()
	logger.InfoContext(ctx, "Shutting down server")
}
