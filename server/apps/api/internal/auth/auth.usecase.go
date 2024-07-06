package auth

import (
	"context"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/user"
	jwt "github.com/boomchanotai/assets-tracker/server/apps/api/internal/utils"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists       = errors.New("EMAIL_ALREADY_EXISTS")
	ErrIncorrectEmailOrPassword = errors.New("INCORRECT_EMAIL_OR_PASSWORD")
)

type usecase struct {
	userRepo  user.Repository
	jwtConfig *jwt.Config
}

func NewUsecase(userRepo user.Repository, jwtConfig *jwt.Config) *usecase {
	return &usecase{
		userRepo:  userRepo,
		jwtConfig: jwtConfig,
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

func (u *usecase) generateAuthToken(ctx context.Context, user entity.User) (accessToken, refreshToken string, exp int64, err error) {
	cachedToken, accessToken, refreshToken, exp, err := jwt.GenerateTokenPair(&user, u.jwtConfig.AccessTokenSecret, u.jwtConfig.RefreshTokenSecret)
	if err != nil {
		return "", "", 0, errors.Wrap(err, "failed to generate token pair")
	}

	if err := u.userRepo.SetUserAuthToken(ctx, user.ID, *cachedToken); err != nil {
		return "", "", 0, errors.Wrap(err, "failed to set user auth token")
	}

	return accessToken, refreshToken, exp, nil
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

	accessToken, refreshToken, exp, err := u.generateAuthToken(ctx, *user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate auth token")
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

	accessToken, refreshToken, exp, err := u.generateAuthToken(ctx, *user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate auth token")
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          exp,
	}, nil
}

func (u *usecase) GetProfile(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := u.userRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by id")
	}

	return user, nil
}

func (u *usecase) RefreshToken(ctx context.Context, token string) (*entity.Token, error) {
	// Refresh Token
	claims, err := jwt.ParseToken(token, u.jwtConfig.RefreshTokenSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse refresh token")
	}

	cachedToken, err := u.userRepo.GetUserAuthToken(ctx, claims.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user auth token")
	}

	if err := jwt.ValidateToken(cachedToken, claims, true); err != nil {
		return nil, errors.Wrap(err, "failed to validate refresh token")
	}

	user, err := u.userRepo.GetUser(ctx, claims.ID)
	if err != nil {
		return nil, errors.Wrap(err, "user not found")
	}

	accessToken, refreshToken, exp, err := u.generateAuthToken(ctx, *user)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate auth token")
	}

	return &entity.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Exp:          exp,
	}, nil
}

func (u *usecase) Logout(ctx context.Context, userID uuid.UUID) error {
	err := u.userRepo.DeleteUserAuthToken(ctx, userID)
	if err != nil {
		return errors.Wrap(err, "failed to delete user auth token")
	}
	return nil
}
