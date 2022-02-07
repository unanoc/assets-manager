package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trustwallet/go-libs/mq"
)

func NewRouter(mq *mq.Client) http.Handler {
	router := gin.Default()

	// middlewares
	SetupMiddlewares(router)

	// metrics
	SetupMetrics(router)

	// routes
	NewValidationAPI().Setup(router)
	NewValuesAPI().Setup(router)
	NewGithubAPI(mq).Setup(router)

	return router
}

func (api *ValidationAPI) Setup(router *gin.Engine) {
	router.POST("/v1/validate/asset_info", api.ValidateAssetInfo)
	router.GET("/v1/validate/url/status", api.CheckURLStatus)
}

func (api *ValuesAPI) Setup(router *gin.Engine) {
	router.GET("/v1/values/tags", api.GetTagValues)
}

func (api *GithubAPI) Setup(router *gin.Engine) {
	router.GET("/v1/github/oauth", api.RedirectToOauth)
	router.GET("/v1/github/oauth/callback", api.HandleOauthCallback)
	router.POST("/v1/github/events/webhook", api.HandleEventsWebhook)
}
