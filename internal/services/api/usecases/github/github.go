package github

import (
	"fmt"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/go-libs/client"
)

type Instance struct {
	client client.Request
}

func New() *Instance {
	return &Instance{
		client: client.InitJSONClient(config.Default.Github.BaseURL, nil),
	}
}

func (i *Instance) GetGithubAccessToken(code string) (string, error) {
	req := &AccessTokenRequest{
		ClientID:     config.Default.Github.ClientID,
		ClientSecret: config.Default.Github.ClientSecret,
		Code:         code,
	}

	var resp AccessTokenResponse
	err := i.client.Post(&resp, "/login/oauth/access_token", req)
	if err != nil {
		return "", fmt.Errorf("failed to get github access token: %w", err)
	}

	return resp.AccessToken, nil
}
