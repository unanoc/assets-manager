package queue

import ghlib "github.com/google/go-github/v38/github"

type EventType string

const (
	PullRequestOpened              = "pull_request_opened"
	PullRequestSynchronize         = "pull_request_synchronize"
	PullRequestReviewCommentOpened = "pull_request_review_comment_opened"
	IssueCommentCreated            = "issue_comment_created"
)

type (
	GithubEventMessage struct {
		Type                     string                               `json:"type"`
		PullRequest              *ghlib.PullRequestEvent              `json:"pull_request"`
		IssueComment             *ghlib.IssueCommentEvent             `json:"issue_comment"`
		PullRequestReviewComment *ghlib.PullRequestReviewCommentEvent `json:"pull_request_review_comment"`
	}
)
