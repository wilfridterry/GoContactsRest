package rest

import (
	domain "contact-list/internal/domain/contact"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Contacts
}

func (h *Handler) InitRouter() *gin.Engine {
	r := gin.Default()

	contactGroup := r.Group("/contacts")
	{
		contactGroup.GET("/", h.getContacts)
	}

	return r
}

type Contacts interface {
	All() ([]domain.Contact, error)
	GetOne(id int64) (domain.Contact, error)
}

func NewHandler(service Contacts) *Handler {
	return &Handler{service}
}

func (h *Handler) getContacts(c *gin.Context) {

} 