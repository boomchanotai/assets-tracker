package model

import (
	"time"

	"github.com/boomchanotai/assets-tracker/server/apps/api/internal/entity"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type User struct {
	ID        uuid.UUID `gorm:"id"`
	Email     string    `gorm:"email"`
	Name      string    `gorm:"name"`
	Password  string    `gorm:"password"`
	CreatedAt time.Time `gorm:"created_at"`
	UpdatedAt time.Time `gorm:"updated_at"`
}

type Account struct {
	ID        uuid.UUID          `gorm:"id"`
	UserID    uuid.UUID          `gorm:"references:User"`
	Type      entity.AccountType `gorm:"type:text"`
	Name      string             `gorm:"name"`
	Bank      string             `gorm:"bank"`
	Balance   decimal.Decimal    `gorm:"balance"`
	CreatedAt time.Time          `gorm:"created_at"`
	UpdatedAt time.Time          `gorm:"updated_at"`
}
