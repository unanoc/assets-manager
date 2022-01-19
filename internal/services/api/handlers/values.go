package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/trustwallet/assets-manager/internal/services/api/usecases/values"
)

type ValuesAPI struct {
	values *values.Instance
}

func NewValuesAPI() API {
	return &ValuesAPI{
		values: values.New(),
	}
}

// @Description Get tag values
// @Router /v1/values/tags [get]
func (api *ValuesAPI) GetTagValues(c *gin.Context) {
	c.JSON(http.StatusOK, api.values.GetTagValues())
}
