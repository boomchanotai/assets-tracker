package user

import (
	"context"
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/interfaces"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/jwt"
	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/model"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
)

const (
	AuthTokenKey = "auth:token"
)

type repository struct {
	db          *gorm.DB
	redisClient *redis.Client
	jwtConfig   *jwt.Config
}

func NewRepository(db *gorm.DB, redisClient *redis.Client, jwtConfig *jwt.Config) interfaces.UserRepository {
	db.AutoMigrate(&model.User{})

	return &repository{
		db:          db,
		redisClient: redisClient,
		jwtConfig:   jwtConfig,
	}
}

func (r *repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, errors.Wrap(err, "can't get users")
	}

	var result []entity.User
	for _, u := range users {
		result = append(result, entity.User{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}

	return result, nil
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {

	var u model.User
	if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, errors.Wrap(err, "can't get user")
	}

	return &entity.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var u model.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, errors.Wrap(err, "can't get user by email")
	}

	return &entity.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}

func (r *repository) CreateUser(ctx context.Context, input entity.UserInput) (*entity.User, error) {
	newUser := model.User{
		ID:       uuid.New(),
		Email:    input.Email,
		Name:     input.Name,
		Password: input.Password,
	}

	if err := r.db.Create(&newUser).Error; err != nil {
		return nil, errors.Wrap(err, "can't create user")
	}

	return &entity.User{
		ID:    newUser.ID,
		Email: newUser.Email,
		Name:  newUser.Name,
	}, nil
}

func (r *repository) UpdateUser(ctx context.Context, id uuid.UUID, input entity.UserInput) (*entity.User, error) {
	panic("not implemented")
}

func getTokenKey(userID uuid.UUID) string {
	return AuthTokenKey + ":" + userID.String()
}

type tokenUID struct {
	AccessUID  uuid.UUID `msgpack:"access_uid"`
	RefreshUID uuid.UUID `msgpack:"refresh_uid"`
}

func (r *repository) SetUserAuthToken(ctx context.Context, userID uuid.UUID, token entity.CachedTokens) error {
	cachedToken, err := msgpack.Marshal(tokenUID{
		AccessUID:  token.AccessUID,
		RefreshUID: token.RefreshUID,
	})
	if err != nil {
		return errors.Wrap(err, "can't marshal cached token")
	}

	err = r.redisClient.Set(ctx, getTokenKey(userID), string(cachedToken), time.Second*time.Duration(r.jwtConfig.AutoLogout)).Err()
	if err != nil {
		return errors.Wrap(err, "can't set token")
	}

	return nil
}

func (r *repository) GetUserAuthToken(ctx context.Context, userID uuid.UUID) (*entity.CachedTokens, error) {
	redisToken, err := r.redisClient.Get(ctx, getTokenKey(userID)).Bytes()
	if err != nil {
		return nil, errors.Wrap(err, "can't get token")
	}

	cachedToken := &tokenUID{}
	err = msgpack.Unmarshal(redisToken, cachedToken)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal cached token")
	}

	return &entity.CachedTokens{
		AccessUID:  cachedToken.AccessUID,
		RefreshUID: cachedToken.RefreshUID,
	}, nil
}

func (r *repository) DeleteUserAuthToken(ctx context.Context, userID uuid.UUID) error {
	err := r.redisClient.Del(ctx, getTokenKey(userID)).Err()
	if err != nil {
		return errors.Wrap(err, "can't delete token")
	}

	return nil
}
