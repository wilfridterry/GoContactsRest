package service

import (
	"contact-list/internal/domain"
	"context"
	"time"
)

type UserRepository interface {
	Create(context.Context, domain.User) error
}

type Hashier interface {
	Hash(string) (string, error)
}

type Users struct {
	repo UserRepository
	hashier Hashier
}

func (service *Users) SignUp(ctx context.Context, inp *domain.UserSignUp) (*domain.User, error) {
	password, err := service.hashier.Hash(inp.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name: inp.Name,
		Email: inp.Password,
		Password: password,
		RegisteredAt: time.Now(),
	}

	return &user, nil
}

func NewUsers(repo UserRepository, hashier Hashier) *Users {
	return &Users{repo, hashier}
}