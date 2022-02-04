package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	ErrorResponse struct {
		Error ErrorDetails `json:"error"`
	}

	ErrorDetails struct {
		Message string `json:"message"`
	}
)

func errorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorDetails{
			Message: message,
		},
	}
}

func abortWithStatusJSON(c *gin.Context, code int) {
	c.AbortWithStatusJSON(code, errorResponse(http.StatusText(code)))
}
