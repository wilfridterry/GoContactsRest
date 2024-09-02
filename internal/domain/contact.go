package domain

import (
	"time"
)

type Contact struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	Address   string    `json:"address"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SaveInputContact struct {
	Name     string `json:"name" binding:"required"`
	LastName string `json:"last_name" binding:"required"`
	Phone    string `json:"phone" binding:"required,e164"`
	Email    string `json:"email" binding:"required,email,unique"`
	Address  string `json:"address" binding:"required"`
	Author   string `json:"author" binding:"required"`
}
