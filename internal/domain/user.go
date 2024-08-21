package domain

import (
	"errors"
	"time"
)

var (
	ErrNotFoundUser = errors.New("Not found user")
)

type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	RegisteredAt time.Time `json:"registered_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SignUpInput struct {
	Name     string `json:"name" binding:"required,gte=2,lte=255"`
	Email    string `json:"email" binding:"required,email,gte=4,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=70"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required,email,gte=4,lte=255"`
	Password string `json:"password" binding:"required,gte=6,lte=70"`
}
