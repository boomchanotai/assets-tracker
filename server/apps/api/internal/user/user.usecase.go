package user

import "github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"

type usecase struct {
	userRepo interfaces.UserRepository
}

func NewUsecase(userRepo interfaces.UserRepository) *usecase {
	return &usecase{
		userRepo: userRepo,
	}
}
