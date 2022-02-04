package handlers

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-libs/mq"
)

func NewRouter(mq *mq.Client) http.Handler {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{config.Default.Clients.AssetsManager.App},
		AllowMethods: []string{"GET", "OPTIONS", "POST", "PUT"},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Content-Length",
			"Accept-Encoding", "Authorization", "Accept", "Cache-Control",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	NewValidationAPI().Setup(router)
	NewValuesAPI().Setup(router)
	NewGithubAPI(mq).Setup(router)

	return router
}

func (api *ValidationAPI) Setup(router *gin.Engine) {
	v1 := router.Group("/v1/")

	{
		v1.POST("/validate/asset_info", api.ValidateAssetInfo)
		v1.GET("/validate/url/status", api.CheckURLStatus)
	}
}

func (api *ValuesAPI) Setup(router *gin.Engine) {
	v1 := router.Group("/v1/")

	{
		v1.GET("/values/tags", api.GetTagValues)
	}
}

func (api *GithubAPI) Setup(router *gin.Engine) {
	v1 := router.Group("/v1/")

	{
		v1.GET("/github/oauth", api.RedirectToOauth)
		v1.GET("/github/oauth/callback", api.HandleOauthCallback)
		v1.POST("/github/events/webhook", api.HandleEventsWebhook)
	}
}
