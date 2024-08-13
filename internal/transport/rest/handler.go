package rest

import (
	domain "contact-list/internal/domain/contact"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Contacts
}

type Contacts interface {
	All(context.Context) ([]domain.Contact, error)
	GetOne(context.Context, int64) (*domain.Contact, error)
	Create(context.Context, *domain.SaveInputContact) error
	Update(context.Context, int64, *domain.SaveInputContact) error
	Delete(context.Context, int64) error
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	v1 := r.Group("/api/v1")
	{
		contacts := v1.Group("/contacts")
		{
			contacts.GET("/", h.getContacts)
			contacts.GET("/:id", h.getContact)
		}
	}

	return r
}


func NewHandler(service Contacts) *Handler {
	return &Handler{service}
}

func (h *Handler) getContacts(c *gin.Context) {
	contacts, err := h.service.All(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, contacts)
	}
} 

func (h *Handler) getContact(c *gin.Context) {
	var id int64
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	contact, err := h.service.GetOne(c.Request.Context(), id)

	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, contact)

} 