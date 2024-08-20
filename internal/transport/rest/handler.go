package rest

import (
	"contact-list/internal/domain"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "contact-list/docs"
)

type Handler struct {
	serviceContacts Contacts
	serviceUsers Users
}

type Contacts interface {
	All(context.Context) ([]domain.Contact, error)
	GetOne(context.Context, int64) (*domain.Contact, error)
	Create(context.Context, *domain.SaveInputContact) error
	Update(context.Context, int64, *domain.SaveInputContact) error
	Delete(context.Context, int64) error
}

type Users interface {
	SignUp(context.Context, *domain.SignUpInput) (*domain.User, error)
	SingIn(context.Context, *domain.SignInInput) (string, error)
}

type Uri struct{
	ID int64 `uri:"id" binding:"required"`
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		contacts := v1.Group("/contacts")
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
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}

func NewHandler(contacts Contacts, users Users) *Handler {
	return &Handler{contacts, users}
}
