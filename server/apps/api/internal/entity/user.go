package entity

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID
	Email     string
	Name      string
	CreatedAt string
	UpdatedAt string
}

func (u User) String() string {
	return u.Name
}
