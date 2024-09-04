package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/wilfridterry/contact-list/internal/domain"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	// audit "github.com/wilfridterry/audit-log/pkg/domain"
)

type UserRepository interface {
	Create(context.Context, *domain.User) (int64, error)
	GetByEmailAndPassword(context.Context, string, string) (*domain.User, error)
}

type SessionRepository interface {
	Create(context.Context, *domain.RefreshSession) error
	GetByToken(context.Context, string) (*domain.RefreshSession, error)
}

type Hashier interface {
	Hash(string) (string, error)
}

type Auth struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	auditClient AuditClient
	auditLog    AuditLog
	hashier     Hashier
	hmacSecret  []byte
	ttlToken    time.Duration
}

type UserClaim struct {
	jwt.RegisteredClaims
	ID        int64
	IssuedAt  int64
	ExpiresAt int64
}

func New(userRepo UserRepository, sessionRepo SessionRepository, auditClient AuditClient, auditLog AuditLog, hashier Hashier, secret []byte, ttlToken time.Duration) *Auth {
	return &Auth{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		auditClient: auditClient,
		auditLog:    auditLog,
		hashier:     hashier,
		hmacSecret:  secret,
		ttlToken:    ttlToken,
	}
}

func (service *Auth) SignUp(ctx context.Context, inp *domain.SignUpInput) (*domain.User, error) {
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

	id, err := service.userRepo.Create(ctx, &user)
	if err != nil {
		return nil, err
	}

	user.ID = id

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action:    audit.ACTION_REGISTER,
	// 	Entity:    audit.ENTITY_USER,
	// 	EntityID:  user.ID,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Users.SignUp",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_REGISTER,
		Entity:    ENTITY_USER,
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignUp",
		}).Error("failed to send log request:", err)
	}

	return &user, nil
}

func (service *Auth) SingIn(ctx context.Context, inp *domain.SignInInput) (string, string, error) {
	password, err := service.hashier.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := service.userRepo.GetByEmailAndPassword(ctx, inp.Email, password)

	if err != nil {
		return "", "", err
	}

	// if err := service.auditClient.SendLogRequest(ctx, audit.LogItem{
	// 	Action:    audit.ACTION_LOGIN,
	// 	Entity:    audit.ENTITY_USER,
	// 	EntityID:  user.ID,
	// 	Timestamp: time.Now(),
	// }); err != nil {
	// 	logrus.WithFields(logrus.Fields{
	// 		"method": "Users.SignIn",
	// 	}).Error("failed to send log request:", err)
	// }

	if err := service.auditLog.Log(LogMessage{
		Action:    ACTION_LOGIN,
		Entity:    ENTITY_USER,
		EntityID:  user.ID,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"method": "Users.SignIn",
		}).Error("failed to send log request:", err)
	}

	return service.generateTokens(ctx, user.ID)
}

func (service *Auth) ParseJWTToken(ctx context.Context, tokenString string) (int64, error) {
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

func (service *Auth) RefreshTokens(ctx context.Context, token string) (string, string, error) {
	session, err := service.sessionRepo.GetByToken(ctx, token)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", domain.ErrRefreshTokenExpired
	}

	return service.generateTokens(ctx, session.UserId)
}

func (service *Auth) generateTokens(ctx context.Context, userId int64) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{},
		ID:               int64(userId),
		IssuedAt:         time.Now().Unix(),
		ExpiresAt:        time.Now().Add(service.ttlToken).Unix(),
	})

	accessToken, err := t.SignedString(service.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := service.newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := service.sessionRepo.Create(ctx, &domain.RefreshSession{
		UserId:    userId,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (service *Auth) newRefreshToken() (string, error) {
	b := make([]byte, 32)
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
