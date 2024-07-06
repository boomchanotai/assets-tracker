package auth

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/user"
	jwt "github.com/boomchanotai/assets-tracker/server/apps/api/internal/utils"
	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists       = errors.New("EMAIL_ALREADY_EXISTS")
	ErrIncorrectEmailOrPassword = errors.New("INCORRECT_EMAIL_OR_PASSWORD")
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

func (u *usecase) Register(ctx context.Context, email, name, password string) (*entity.Token, error) {
	if _, err := u.userRepo.GetUserByEmail(ctx, email); err == nil {
		return nil, errors.Wrap(ErrEmailAlreadyExists, "email already exists")
	}

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

	// TODO: SECRET should be stored in config
	cachedToken, accessToken, refreshToken, exp, err := jwt.GenerateTokenPair(user, "SECRET", "SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token pair")
	}

	if err := u.userRepo.SetUserAuthToken(ctx, user.ID, *cachedToken); err != nil {
		return nil, errors.Wrap(err, "failed to set user auth token")
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          exp,
	}, nil
}

func (u *usecase) Login(ctx context.Context, email, password string) (*entity.Token, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}

	if !checkPassword(user.Password, password) {
		return nil, errors.Wrap(ErrIncorrectEmailOrPassword, "incorrect email or password")
	}

	// TODO: SECRET should be stored in config
	cachedToken, accessToken, refreshToken, exp, err := jwt.GenerateTokenPair(user, "SECRET", "SECRET")
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token pair")
	}

	if err := u.userRepo.SetUserAuthToken(ctx, user.ID, *cachedToken); err != nil {
		return nil, errors.Wrap(err, "failed to set user auth token")
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          exp,
	}, nil
}
