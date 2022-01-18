package validation

import "github.com/trustwallet/assets-manager/internal/services/api/models"

type (
	AssetInfoRequest struct {
		ID          string   `json:"id,omitempty"`
		Name        string   `json:"name,omitempty"`
		Type        string   `json:"type,omitempty"`
		Symbol      string   `json:"symbol,omitempty"`
		Website     string   `json:"website,omitempty"`
		Explorer    string   `json:"explorer,omitempty"`
		Description string   `json:"description,omitempty"`
		Status      string   `json:"status,omitempty"`
		Decimals    int      `json:"decimals,omitempty"`
		Links       []Link   `json:"links,omitempty"`
		Tags        []string `json:"tags,omitempty"`
	}

	Link struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	}

	AssetInfoResponse struct {
		Status models.StatusType `json:"status"`
		Errors []Error           `json:"errors"`
	}

	Error struct {
		Message string `json:"message"`
	}
)
