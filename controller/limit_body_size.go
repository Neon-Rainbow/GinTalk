package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LimitBodySizeMiddleware(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		c.Next()
	}
}
