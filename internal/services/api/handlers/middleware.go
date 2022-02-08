package handlers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupMiddlewares(router *gin.Engine) {
	router.Use(GetCORSMiddleware())
}

func GetCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "OPTIONS", "POST", "PUT"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "Authorization", "Accept", "Cache-Control",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
