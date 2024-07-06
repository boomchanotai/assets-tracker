package entity

import "github.com/google/uuid"

type User struct {
	ID    uuid.UUID
	Email string
	Name  string
}

func (u User) String() string {
	return u.Name
}

type UserInput struct {
	Email    string
	Name     string
	Password string
}
