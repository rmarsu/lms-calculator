package http

import (
	"github.com/gin-gonic/gin"
)

func newErrorResponse(code int, msg string, c *gin.Context) {
	c.AbortWithStatusJSON(code, gin.H{
		"error": msg,
	})
}
