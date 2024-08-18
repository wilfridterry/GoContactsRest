package rest

import (
	"contact-list/internal/domain"

	"github.com/swaggo/swag/example/celler/httputil"
	"net/http"
	"github.com/gin-gonic/gin"

)

// SignUp godoc
// @Summary      sign up to the system
// @Description  sign up with data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user body domain.UserSignUp true "User "
// @Success      201
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) { 

	var inp domain.UserSignUp

	if err := c.ShouldBindJSON(&inp); err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)

		return
	}

	user, err := h.serviceUsers.SignUp(c.Request.Context(), &inp)
	
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created.", "user": user})
}