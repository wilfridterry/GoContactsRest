package rest

import (
	domain "contact-list/internal/domain/contact"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Contacts
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(Logger())

	contactGroup := r.Group("/contacts")
	{
		contactGroup.GET("/", h.getContacts)
	}

	return r
}

type Contacts interface {
	All() ([]domain.Contact, error)
	GetOne(id int64) (*domain.Contact, error)
}

func NewHandler(service Contacts) *Handler {
	return &Handler{service}
}

func (h *Handler) getContacts(c *gin.Context) {
	data := map[string]string{
		"hello": "world",
	}
	c.JSON(http.StatusOK, data)
} 