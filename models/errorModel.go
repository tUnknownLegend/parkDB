package models

import (
	"github.com/gin-gonic/gin"
)

type Error struct {
	Message string `json:"message"`
}

func GetErrorResponse(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, Error{Message: message})
}
