package validation

import (
	"fmt"
	"net/http"
	"strings"
)

func (i *Instance) CheckURLStatus(url string) (*URLStatusResponse, error) {
	if strings.Contains(url, "?") {
		return nil, fmt.Errorf("forbidden url: url contains query parameters")
	}

	resp, err := http.Get(url) // nolint:gosec // no need, status code check only
	if err != nil {
		return nil, fmt.Errorf("failed to make get request to %s: %w", url, err)
	}
	defer resp.Body.Close()

	return &URLStatusResponse{
		URL:        url,
		StatusCode: resp.StatusCode,
	}, nil
}
