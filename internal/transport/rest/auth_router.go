package rest

import (
	"contact-list/internal/domain"
	"errors"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/swag/example/celler/httputil"
)

// SignUp godoc
// @Summary      sign up to the system
// @Description  sign up with data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user body domain.UserSignUp true "user sign up"
// @Success      201
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var inp domain.SignUpInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)

		return
	}

	user, err := h.authServie.SignUp(c.Request.Context(), &inp)

	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created.", "user": user})
}

// SignIn godoc
// @Summary      sign in to the system
// @Description  sign in with data
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param user body domain.UserSignIn true "user sign in"
// @Success      200
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /auth/sign-in [get]
func (h *Handler) signIn(c *gin.Context) {
	var inp domain.SignInInput

	if err := c.ShouldBindJSON(&inp); err != nil {
		httputil.NewError(c, http.StatusUnprocessableEntity, err)

		return
	}

	accessToken, refreshToken, err := h.authServie.SingIn(c.Request.Context(), &inp)

	if err != nil {
		if errors.Is(err, domain.ErrNotFoundUser) {
			httputil.NewError(c, http.StatusBadRequest, err)
			return
		}

		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("refresh-token", refreshToken, 3600, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}

// Refresh godoc
// @Summary      refresh tokens
// @Description  refresh tokens for auth
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /auth/sign-in [Get]
func (h *Handler) refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh-token")
	if err != nil {
		httputil.NewError(c, http.StatusBadRequest, err)
		return
	}

	logrus.Infof("%s", cookie)

	accessToken, refreshToken, err := h.authServie.RefreshTokens(c.Request.Context(), cookie)
	if err != nil {
		httputil.NewError(c, http.StatusInternalServerError, err)
		return
	}

	c.SetCookie("refresh-token", refreshToken, 3600, "/", "localhost", true, true)
	c.JSON(http.StatusOK, gin.H{"token": accessToken})
}
