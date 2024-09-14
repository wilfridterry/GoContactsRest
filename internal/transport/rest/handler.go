package rest

import (
	"context"

	"github.com/wilfridterry/contact-list/internal/domain"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/wilfridterry/contact-list/docs"
)

//go:generate mockgen -source=handler.go -destination=mocks/mock.go

type Handler struct {
	contactService Contacts
	authServie     Auth
}

type Contacts interface {
	All(context.Context) ([]domain.Contact, error)
	GetOne(context.Context, int64) (*domain.Contact, error)
	Create(context.Context, *domain.SaveInputContact) error
	Update(context.Context, int64, *domain.SaveInputContact) error
	Delete(context.Context, int64) error
}

type Auth interface {
	SignUp(context.Context, *domain.SignUpInput) (*domain.User, error)
	SingIn(context.Context, *domain.SignInInput) (string, string, error)
	ParseJWTToken(context.Context, string) (int64, error)
	RefreshTokens(context.Context, string) (string, string, error)
}

type Uri struct {
	ID int64 `uri:"id" binding:"required"`
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		contacts := v1.Group("/contacts").Use(h.AuthJWT())
		{
			contacts.POST("/", h.createContact)
			contacts.GET("/", h.getContacts)
			contacts.GET("/:id", h.getContact)
			contacts.DELETE("/:id", h.deleteContact)
			contacts.PUT("/:id", h.updateAccount)
		}

		auth := v1.Group("/auth")
		{
			auth.POST("/sign-up", h.signUp)
			auth.GET("/sign-in", h.signIn)
			auth.GET("/refresh", h.refresh)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func NewHandler(contacts Contacts, auth Auth) *Handler {
	return &Handler{contacts, auth}
}
