package rest

import (
	domain "contact-list/internal/domain/contact"
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag/example/celler/httputil"

	_ "contact-list/docs"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}


func NewHandler(service Contacts) *Handler {
	return &Handler{service}
}

// ListContacts godoc
// @Summary      List contacts
// @Description  get contacts
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Success      200  {array}   domain.Contact
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /contacts [get]
func (h *Handler) getContacts(c *gin.Context) {
	contacts, err := h.service.All(c.Request.Context())
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
	} else {
		c.JSON(http.StatusOK, contacts)
	}
} 

// ShowContact godoc
// @Summary      Show a contact
// @Description  get string by ID
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Contact ID"
// @Success      200  {object}  domain.Contact
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /contacts/{id} [get]
func (h *Handler) getContact(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	contact, err := h.service.GetOne(c.Request.Context(), uri.ID)

	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			httputil.NewError(c, http.StatusNotFound, err)

			return
		}
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, contact)	
} 

// CreateContact godoc
// @Summary      Create a contact
// @Description  create a contact with data
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param contact body domain.SaveInputContact true "Contact paylaod"
// @Success      201
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /contacts [post]
func (h *Handler) createContact(c *gin.Context) {
	var inp domain.SaveInputContact
	if err := c.ShouldBindJSON(&inp); err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)

		return
	}

	if err := h.service.Create(c.Request.Context(), &inp); err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created."})
}


// DeleteContact godoc
// @Summary      Delete a contact
// @Description  delete a contact by ID
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Contact ID"
// @Success      204
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /contacts/{id} [delete]
func (h *Handler) deleteContact(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)

		return
	}

	if err := h.service.Delete(c.Request.Context(), uri.ID); 
	err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusNoContent, gin.H{})	
} 

// UpdateContact godoc
// @Summary      Update a contact
// @Description  Update a contact with data by ID
// @Tags         contacts
// @Accept       json
// @Produce      json
// @Param contact body domain.SaveInputContact true "Contact paylaod"
// @Param        id   path      int  true  "Contact ID"
// @Success      200  {object}  domain.Contact
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /contacts/{id} [put]
func (h *Handler) updateAccount(c *gin.Context) {
	var uri Uri
	if err := c.ShouldBindUri(&uri); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	var inp domain.SaveInputContact
	if err := c.ShouldBindJSON(&inp); err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)

		return
	}

	if err := h.service.Update(c.Request.Context(), uri.ID, &inp); 
	err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	contact, err := h.service.GetOne(c.Request.Context(), uri.ID)
	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			httputil.NewError(c, http.StatusNotFound, err)
			
			return
		}
		httputil.NewError(c, http.StatusBadRequest, err)

		return
	}

	c.JSON(http.StatusOK, contact)
} 