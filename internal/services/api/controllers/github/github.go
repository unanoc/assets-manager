package github

import (
	"encoding/json"
	"fmt"

	ghlib "github.com/google/go-github/v38/github"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/queue"
	"github.com/trustwallet/go-libs/client"
	"github.com/trustwallet/go-libs/mq"
)

const (
	eventActionCreated     = "created"
	eventActionOpened      = "opened"
	eventActionReopened    = "reopened"
	eventActionSynchronize = "synchronize"
)

type Controller struct {
	client client.Request
	queue  mq.Queue
}

func NewController(mq *mq.Client) *Controller {
	return &Controller{
		client: client.InitJSONClient(config.Default.Github.BaseURL, nil),
		queue:  mq.InitQueue(queue.QueueAssetManagerProcessGithubEvent),
	}
}

func (i *Controller) GetGithubAccessToken(code string) (string, error) {
	req := &AccessTokenRequest{
		ClientID:     config.Default.Github.ClientID,
		ClientSecret: config.Default.Github.ClientSecret,
		Code:         code,
	}

	var resp AccessTokenResponse
	if err := i.client.Post(&resp, "/login/oauth/access_token", req); err != nil {
		return "", fmt.Errorf("failed to get github access token: %w", err)
	}

	return resp.AccessToken, nil
}

func (i *Controller) PushGithubEventToQueue(eventPayload interface{}) error {
	switch event := eventPayload.(type) {
	case *ghlib.PullRequestEvent:
		if *event.Action == eventActionOpened || *event.Action == eventActionReopened {
			return publishGithubEvent(queue.GithubEventMessage{
				Type:        queue.PullRequestOpened,
				PullRequest: event,
			}, i.queue)
		}

		if *event.Action == eventActionSynchronize {
			return publishGithubEvent(queue.GithubEventMessage{
				Type:        queue.PullRequestSynchronize,
				PullRequest: event,
			}, i.queue)
		}

	case *ghlib.IssueCommentEvent:
		if *event.Action == eventActionCreated {
			return publishGithubEvent(queue.GithubEventMessage{
				Type:         queue.IssueCommentCreated,
				IssueComment: event,
			}, i.queue)
		}

	case *ghlib.PullRequestReviewCommentEvent:
		if *event.Action == eventActionOpened {
			return publishGithubEvent(queue.GithubEventMessage{
				Type:                     queue.PullRequestReviewCommentOpened,
				PullRequestReviewComment: event,
			}, i.queue)
		}
	}

	return nil
}

func publishGithubEvent(event queue.GithubEventMessage, q mq.Queue) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err = q.Publish(body); err != nil {
		return fmt.Errorf("failed to publish to queue '%s': %w", q.Name(), err)
	}

	return nil
}
