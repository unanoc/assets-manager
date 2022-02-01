package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/services/api/usecases/validation"
)

type ValidationAPI struct {
	validator *validation.Instance
}

func NewValidationAPI() API {
	return &ValidationAPI{
		validator: validation.New(),
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
		log.Error(fmt.Errorf("url not found"))
		abortWithStatusJSON(c, http.StatusBadRequest)

		return
	}

	resp, err := api.validator.CheckURLStatus(url)
	if err != nil {
		handleResponse(c, err)

		return
	}

	c.JSON(200, resp)
}
