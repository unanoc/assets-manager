package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services/api/usecases/github"
)

type GithubAPI struct {
	github *github.Instance
}

func NewGithubAPI() API {
	return &GithubAPI{
		github: github.New(),
	}
}

// @Description Redirects to github oauth
// @Router /v1/github/oauth [get]
func (api *GithubAPI) RedirectToOauth(c *gin.Context) {
	scope := "public_repo%20read:user"

	githubURL := fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&scope=%s",
		config.Default.Github.BaseURL, config.Default.Github.ClientID, scope)

	c.Redirect(http.StatusFound, githubURL)
}

// @Description Gets github access token and redirects to the root of app
// @Router /v1/github/oauth/callback [get]
func (api *GithubAPI) HandleOauthCallback(c *gin.Context) {
	code, ok := c.GetQuery("code")
	if !ok {
		log.Debug(fmt.Errorf("code not found"))
		abortWithStatusJSON(c, http.StatusBadRequest)

		return
	}

	accessToken, err := api.github.GetGithubAccessToken(code)
	if err != nil {
		log.Error(err)
		handleResponse(c, err)

		return
	}

	log.WithField("token", accessToken).Debugf("Received github access token")

	url := fmt.Sprintf("%s/?ghtoken=%s",
		config.Default.Clients.AssetsManager.App, accessToken)

	c.Redirect(http.StatusFound, url)
}
