package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag/example/celler/httputil"
)

type CtxValue int
const (
	ctxUserId CtxValue = iota
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		log.Printf("[%s] %s\n", ctx.Request.Method, ctx.Request.URL)
	}
}

func (h *Handler) AuthJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := getBearerToken(ctx)
		if err != nil {
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		userId, err := h.authServie.ParseJWTToken(ctx.Request.Context(), token)
		if err != nil {
			httputil.NewError(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		rCtx := context.WithValue(ctx.Request.Context(), ctxUserId, userId) 
		ctx.Request = ctx.Request.WithContext(rCtx)
		ctx.Next()
	}
}

func getBearerToken(ctx *gin.Context) (string, error) {
	header := ctx.GetHeader("Authorization")

	if header == "" {
		return "", errors.New("invalid auth header")
	}

	headerSectors := strings.Split(header, " ")

	if len(headerSectors) != 2 && headerSectors[0] == "Bearer" {
		return "", errors.New("invalid auth header")
	}

	return headerSectors[1], nil
}