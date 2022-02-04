package values

import (
	"github.com/trustwallet/assets-manager/internal/config"
)

func (i *Controller) GetTagValues() *TagValuesResponse {
	cfgTags := config.Default.Tags

	tags := make([]Tag, len(cfgTags))
	for i, tag := range cfgTags {
		tags[i] = Tag{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
		}
	}

	return &TagValuesResponse{
		Tags: tags,
	}
}
