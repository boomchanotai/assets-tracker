package user

import (
	"context"
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	GetUsers(ctx context.Context) ([]entity.User, error)
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, user entity.User) error
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
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	db.AutoMigrate(&user{})

	return &repository{
		db: db,
	}
}

func (r *repository) GetUsers(ctx context.Context) ([]entity.User, error) {
	panic("not implemented")
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	panic("not implemented")
}

func (r *repository) CreateUser(ctx context.Context, user entity.User) error {
	panic("not implemented")
}

func (r *repository) UpdateUser(ctx context.Context, user entity.User) error {
	panic("not implemented")
}
