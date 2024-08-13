package domain

import (
	"errors"
	"time"
)

var (
	ErrContactNotFound = errors.New("Contact not found")
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
	Name      *string
	LastName  *string
	Phone     *string
	Email     *string
	Address   *string
	Author    *string
}