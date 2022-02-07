package handlers

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/trustwallet/assets-manager/internal/config"
)

func SetupMiddlewares(router *gin.Engine) {
	router.Use(GetCORSMiddleware())
}

func GetCORSMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins: []string{config.Default.Clients.AssetsManager.App},
		AllowMethods: []string{"GET", "OPTIONS", "POST", "PUT"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "Authorization", "Accept", "Cache-Control",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
