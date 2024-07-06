package user

import (
	"context"
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
	"gorm.io/gorm"
)

const (
	AuthTokenKey = "auth:token"
)

type Repository interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	CreateUser(ctx context.Context, userInput entity.UserInput) (*entity.User, error)
	UpdateUser(ctx context.Context, user entity.User) error

	SetUserAuthToken(ctx context.Context, userID uuid.UUID, token entity.CachedTokens) error
	GetUserAuthToken(ctx context.Context, userID uuid.UUID) (*entity.CachedTokens, error)
	DeleteUserAuthToken(ctx context.Context, userID uuid.UUID) error
}

type user struct {
	ID        uuid.UUID `gorm:"id"`
	Email     string    `gorm:"email"`
	Name      string    `gorm:"name"`
	Password  string    `gorm:"password"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type repository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewRepository(db *gorm.DB, redisClient *redis.Client) Repository {
	db.AutoMigrate(&user{})

	return &repository{
		db:          db,
		redisClient: redisClient,
	}
}

func (r *repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	var users []user
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

	var u user
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
	var u user
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

func (r *repository) CreateUser(ctx context.Context, userInput entity.UserInput) (*entity.User, error) {
	newUser := user{
		ID:        uuid.New(),
		Email:     userInput.Email,
		Name:      userInput.Name,
		Password:  userInput.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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

func (r *repository) UpdateUser(ctx context.Context, user entity.User) error {
	panic("not implemented")
}

func getTokenKey(userID uuid.UUID) string {
	return AuthTokenKey + ":" + userID.String()
}

func (r *repository) SetUserAuthToken(ctx context.Context, userID uuid.UUID, token entity.CachedTokens) error {
	cachedToken, err := msgpack.Marshal(token)
	if err != nil {
		return errors.Wrap(err, "can't marshal cached token")
	}

	// TODO: duration should be stored in config (30 days)
	err = r.redisClient.Set(ctx, getTokenKey(userID), string(cachedToken), time.Minute*60*24*30).Err()
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

	cachedToken := &entity.CachedTokens{}
	err = msgpack.Unmarshal(redisToken, cachedToken)
	if err != nil {
		return nil, errors.Wrap(err, "can't unmarshal cached token")
	}

	return cachedToken, nil
}

func (r *repository) DeleteUserAuthToken(ctx context.Context, userID uuid.UUID) error {
	err := r.redisClient.Del(ctx, getTokenKey(userID)).Err()
	if err != nil {
		return errors.Wrap(err, "can't delete token")
	}

	return nil
}
