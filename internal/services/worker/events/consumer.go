package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/trustwallet/assets-manager/internal/queue"
	"github.com/trustwallet/go-libs/mq"
)

func GetEventConsumer(ctx context.Context, eh *Handler) func(mq.Message) error {
	return func(message mq.Message) error {
		var event queue.GithubEventMessage

		err := json.Unmarshal(message, &event)
		if err != nil {
			return fmt.Errorf("failed to unmarshal Github event: %w", err)
		}

		switch event.Type {
		case queue.PullRequestOpened:
			err = eh.HandlePullRequestOpened(ctx, event.PullRequest)
		case queue.PullRequestSynchronize:
			err = eh.HandlePullRequestChangesPushed(ctx, event.PullRequest)
		case queue.IssueCommentCreated:
			err = eh.HandleIssueCommentCreated(ctx, event.IssueComment)
		case queue.PullRequestReviewCommentOpened:
			err = eh.HandlePullRequestReviewCommentCreated(ctx, event.PullRequestReviewComment)
		}

		if err != nil {
			return fmt.Errorf("failed to handle Github event: %w", err)
		}

		return nil
	}
}
