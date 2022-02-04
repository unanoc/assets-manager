package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	ghlib "github.com/google/go-github/v38/github"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services/api/controllers/github"
	"github.com/trustwallet/go-libs/mq"
)

type GithubAPI struct {
	github *github.Controller
}

func NewGithubAPI(mq *mq.Client) API {
	return &GithubAPI{
		github: github.NewController(mq),
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
		log.WithError(err).Error("github access token getting error")
		abortWithStatusJSON(c, http.StatusInternalServerError)

		return
	}

	url := fmt.Sprintf("%s/?ghtoken=%s",
		config.Default.Clients.AssetsManager.App, accessToken)

	c.Redirect(http.StatusFound, url)
}

// @Description Gets github events and pushes them to queue.
// @Router /v1/github/events/webhook [post]
func (api *GithubAPI) HandleEventsWebhook(c *gin.Context) {
	// Validate webhook payload.
	payloadBytes, err := ghlib.ValidatePayload(c.Request, []byte(config.Default.Github.AppWebhookSecret))
	if err != nil {
		log.WithError(err).Error("github payload validation error")
		abortWithStatusJSON(c, http.StatusForbidden)

		return
	}

	// Parse the incoming request into an event.
	eventPayload, err := ghlib.ParseWebHook(ghlib.WebHookType(c.Request), payloadBytes)
	if err != nil {
		log.WithError(err).Error("parse github web hook error")
		abortWithStatusJSON(c, http.StatusBadRequest)

		return
	}

	log.WithField("event", ghlib.WebHookType(c.Request)).Debug("Incoming event")

	err = api.github.PushGithubEventToQueue(eventPayload)
	if err != nil {
		log.WithError(err).WithField("event", ghlib.WebHookType(c.Request)).Error("push github event to queue error")
		abortWithStatusJSON(c, http.StatusInternalServerError)
	}

	c.JSON(http.StatusOK, nil)
}
