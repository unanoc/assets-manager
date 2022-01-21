package handlers

import (
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

	response, err := api.validator.ValidateAssetInfo(request)
	if err != nil {
		log.Error(err)
		handleResponse(c, err)

		return
	}

	c.JSON(http.StatusOK, response)
}
