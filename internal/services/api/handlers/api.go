package handlers

import "github.com/gin-gonic/gin"

// API contains generic methods that should be used by other API interfaces.
type API interface {
	Setup(router *gin.Engine)
}
