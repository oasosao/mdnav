package middleware

import (
	"net/http"

	"mdnav/internal/core"

	"github.com/gin-gonic/gin"
)

func Options(ctx *core.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Writer.Header().Set("Access-Control-Max-Age", "86400")
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}
