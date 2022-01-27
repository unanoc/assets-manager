package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() http.Handler {
	router := gin.Default()

	NewValidationAPI().Setup(router)
	NewValuesAPI().Setup(router)
	NewGithubAPI().Setup(router)

	return router
}

func (api *ValidationAPI) Setup(router *gin.Engine) {
	v1 := router.Group("/v1/")

	{
		v1.POST("/validate/asset_info", api.ValidateAssetInfo)
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
		v1.GET("/github/callback", api.HandleCallback)
	}
}
