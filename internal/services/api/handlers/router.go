package handlers

import "github.com/gin-gonic/gin"

func NewRouter() *gin.Engine {
	router := gin.Default()

	validationAPI := NewValidationAPI()
	validationAPI.Setup(router)

	return router
}

func (api *ValidationAPI) Setup(router *gin.Engine) {
	v1 := router.Group("/v1/")
	{
		v1.POST("/validate/asset_info", api.ValidateAssetInfo)
	}
}
