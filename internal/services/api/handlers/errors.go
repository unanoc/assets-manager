package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

// Generic implementation for handling response.
func handleResponse(c *gin.Context, err error) {
	switch errors.Is { // nolint:gocritic // it shows the way for handling new errors
	default:
		log.Debug(err)
		abortWithStatusJSON(c, http.StatusInternalServerError)
	}
}
