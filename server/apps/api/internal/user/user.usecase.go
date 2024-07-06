package user

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
)

type usecase struct {
	userRepo Repository
}

func NewUsecase(userRepo Repository) *usecase {
	return &usecase{
		userRepo: userRepo,
	}
}

func (u *usecase) Login(ctx context.Context, email, password string) (*entity.User, error) {
	panic("not implemented")
}
