package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) String() string {
	return u.Name
}

type UserInput struct {
	Email    string
	Name     string
	Password string
}

type Token struct {
	AccessToken  string
	RefreshToken string
	Exp          int64
}

type CachedTokens struct {
	AccessUID  uuid.UUID `json:"access"`
	RefreshUID uuid.UUID `json:"refresh"`
}
