package handler

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponseAPI struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Payload    interface{} `json:"payload,omitempty"`
}

func Success(c *gin.Context, status int, message string, data gin.H) {
	res := SuccessResponseAPI{
		StatusCode: status,
		Message:    message,
		Payload:    data,
	}
	c.JSON(status, res)
}
