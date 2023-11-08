package handler

import (
	"github.com/gin-gonic/gin"
)

type ErrorResponseAPI struct {
	StatusCode int    `json:"status_code"`
	ErrorCode  int    `json:"error_code"`
	Message    string `json:"message"`
}

func Error(c *gin.Context, status int, errorCode int, message string) {
	res := ErrorResponseAPI{
		StatusCode: status,
		ErrorCode:  errorCode,
		Message:    message,
	}
	c.AbortWithStatusJSON(status, res)
}
