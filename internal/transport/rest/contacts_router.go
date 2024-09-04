package rest

import (
	"github.com/wilfridterry/contact-list/internal/domain"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

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
	contacts, err := h.contactService.All(c.Request.Context())
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

	contact, err := h.contactService.GetOne(c.Request.Context(), uri.ID)

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

	if err := h.contactService.Create(c.Request.Context(), &inp); err != nil {
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

	if err := h.contactService.Delete(c.Request.Context(), uri.ID); err != nil {
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

	if err := h.contactService.Update(c.Request.Context(), uri.ID, &inp); err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	contact, err := h.contactService.GetOne(c.Request.Context(), uri.ID)
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
