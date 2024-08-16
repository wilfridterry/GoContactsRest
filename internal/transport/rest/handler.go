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
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.service.GetOne(c.Request.Context(), uri.ID)

	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contact)	
} 

func (h *Handler) createContact(c *gin.Context) {
	var inp domain.SaveInputContact
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(c.Request.Context(), &inp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created."})
}

func (h *Handler) deleteContact(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Delete(c.Request.Context(), uri.ID); 
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})	
} 

func (h *Handler) updateAccount(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var inp domain.SaveInputContact
	if err := c.ShouldBindJSON(&inp); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(c.Request.Context(), uri.ID, &inp); 
	err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contact, err := h.service.GetOne(c.Request.Context(), uri.ID)
	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, contact)
} 