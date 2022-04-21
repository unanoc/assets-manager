package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/services/api/controllers/validation"
)

type ValidationAPI struct {
	validator *validation.Controller
}

func NewValidationAPI() API {
	return &ValidationAPI{
		validator: validation.NewController(),
	}
}

// @Description Validate asset info
// @Router /v1/validate/asset_info [post]
func (api *ValidationAPI) ValidateAssetInfo(c *gin.Context) {
	var request validation.AssetInfoRequest

	if err := c.Bind(&request); err != nil {
		abortWithStatusJSON(c, http.StatusBadRequest)

		return
	}

	response := api.validator.ValidateAssetInfo(request)

	c.JSON(http.StatusOK, response)
}

// @Description Checks URL's status
// @Router /v1/validate/url/status [get]
func (api *ValidationAPI) CheckURLStatus(c *gin.Context) {
	url, ok := c.GetQuery("url")
	if !ok {
		log.Debug(fmt.Errorf("url not found"))
		abortWithStatusJSON(c, http.StatusBadRequest)

		return
	}

	resp, err := api.validator.CheckURLStatus(url)
	if err != nil {
		abortWithStatusJSON(c, http.StatusInternalServerError)

		return
	}

	c.JSON(200, resp)
}
