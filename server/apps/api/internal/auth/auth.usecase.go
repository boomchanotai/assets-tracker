package auth

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/user"
	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/bcrypt"
)

type usecase struct {
	userRepo user.Repository
}

func NewUsecase(userRepo user.Repository) *usecase {
	return &usecase{
		userRepo: userRepo,
	}
}

func getHashPassword(password string) (string, error) {
	bytePassword := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func checkPassword(hashPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func (u *usecase) Register(ctx context.Context, email, name, password string) (*entity.User, error) {
	hashPassword, err := getHashPassword(password)
	if err != nil {
		return nil, errors.Wrap(err, "failed to hash password")
	}

	user, err := u.userRepo.CreateUser(ctx, entity.UserInput{
		Email:    email,
		Name:     name,
		Password: hashPassword,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}

	return user, nil
}

func (u *usecase) Login(email, password string) error {
	panic("not implemented")
}
