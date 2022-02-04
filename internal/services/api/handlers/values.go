package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trustwallet/assets-manager/internal/services/api/controllers/values"
)

type ValuesAPI struct {
	values *values.Controller
}

func NewValuesAPI() API {
	return &ValuesAPI{
		values: values.NewController(),
	}
}

// @Description Get tag values
// @Router /v1/values/tags [get]
func (api *ValuesAPI) GetTagValues(c *gin.Context) {
	c.JSON(http.StatusOK, api.values.GetTagValues())
}
