package service

import (
	"contact-list/internal/domain"
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserRepository interface {
	Create(context.Context, domain.User) (int64, error)
	GetByEmailAndPassword(context.Context, string, string) (*domain.User, error)
}

type Hashier interface {
	Hash(string) (string, error)
}

type Users struct {
	repo       UserRepository
	hashier    Hashier
	hmacSecret []byte
	ttlToken   time.Duration
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID        int64
	IssuedAt  int64
	ExpiresAt int64
}

func NewUsers(repo UserRepository, hashier Hashier, secret []byte, ttlToken time.Duration) *Users {
	return &Users{repo, hashier, secret, ttlToken}
}

func (service *Users) SignUp(ctx context.Context, inp *domain.SignUpInput) (*domain.User, error) {
	password, err := service.hashier.Hash(inp.Password)
	if err != nil {
		return nil, err
	}

	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	id, err := service.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	return &user, nil
}

func (service *Users) SingIn(ctx context.Context, inp *domain.SignInInput) (string, error) {
	password, err := service.hashier.Hash(inp.Password)
	if err != nil {
		return "", err
	}

	user, err := service.repo.GetByEmailAndPassword(ctx, inp.Email, password)

	if err != nil {
		return "", err
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		ID:               int64(user.ID),
		IssuedAt:         time.Now().Unix(),
		ExpiresAt:        time.Now().Add(service.ttlToken).Unix(),
	})

	return t.SignedString(service.hmacSecret)
}

func (service *Users) ParseJWTToken(ctx context.Context, tokenString string) (int64, error) {
	userClaim := &UserClaim{}
	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return service.hmacSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, err
	}

	userClaim, ok := token.Claims.(*UserClaim)
	
	if !ok {
		return 0, err
	}

	return userClaim.ID, nil
}
