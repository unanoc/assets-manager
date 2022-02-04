package validation

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (i *Controller) CheckURLStatus(websiteURL string) (*URLStatusResponse, error) {
	if strings.Contains(websiteURL, "?") {
		return &URLStatusResponse{
			URL:           websiteURL,
			StatusCode:    http.StatusBadRequest,
			StatusMessage: "forbidden url: url contains query parameters",
		}, nil
	}

	resp, err := http.Get(websiteURL) // nolint:gosec // no need, status code check only
	if err != nil {
		var err2 *url.Error
		if errors.As(err, &err2) {
			return &URLStatusResponse{
				URL:           websiteURL,
				StatusCode:    http.StatusNotFound,
				StatusMessage: fmt.Sprintf("failed to make get request to %s: %s", websiteURL, err2),
			}, nil
		}

		return nil, err
	}
	defer resp.Body.Close()

	return &URLStatusResponse{
		URL:           websiteURL,
		StatusCode:    resp.StatusCode,
		StatusMessage: resp.Status,
	}, nil
}
