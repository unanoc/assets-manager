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

func (api *GithubAPI) RedirectToOauth(c *gin.Context) {
	scope := "public_repo%20read:user"

	githubURL := fmt.Sprintf("%s/login/oauth/authorize?client_id=%s&scope=%s",
		config.Default.Github.BaseURL, config.Default.Github.ClientID, scope)

	c.Redirect(http.StatusFound, githubURL)
}

func (api *GithubAPI) HandleCallback(c *gin.Context) {
	code, ok := c.GetQuery("code")
	if !ok {
		log.Error(fmt.Errorf("code not found"))
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
