package rest

import (
	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		log.Printf("[%s] %s\n", ctx.Request.Method, ctx.Request.URL)
	}
}