package github

import (
	"context"
	"net/http"

	ghi "github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v38/github"
	"github.com/pkg/errors"

	"github.com/trustwallet/assets-manager/internal/merge-fee-bot/config"
)

type Client struct {
	client *github.Client
}

// NewClient return an instance of Github for working with Github API.
func NewClient() (*Client, error) {
	tr, err := ghi.NewAppsTransport(http.DefaultTransport, config.Default.Github.AppID,
		[]byte(config.Default.Github.AppPrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a transport without installation")
	}

	client, err := github.NewEnterpriseClient(
		config.Default.Github.BaseURL, config.Default.Github.BaseURL, &http.Client{Transport: tr})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a github client")
	}

	insts, _, err := client.Apps.ListInstallations(context.Background(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get list of installations")
	}

	tr2, err := ghi.New(
		http.DefaultTransport, config.Default.Github.AppID, insts[0].GetID(), []byte(config.Default.Github.AppPrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a transport with installation")
	}

	client, err = github.NewEnterpriseClient(
		config.Default.Github.BaseURL, config.Default.Github.BaseURL, &http.Client{Transport: tr2})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create a github client")
	}

	return &Client{client: client}, nil
}

// SetLabelOnPullRequqest sets a new label on a pull request if label does not exist.
func (c *Client) SetLabelOnPullRequest(ctx context.Context, owner, repo string, prNum int, label *github.Label) error {
	allLabels, _, err := c.client.Issues.ListLabels(ctx, owner, repo, nil)
	if err != nil {
		return errors.Wrap(err, "failed to get labels list")
	}

	var labelAlreadyExist bool

	for _, l := range allLabels {
		if l.GetName() == label.GetName() {
			labelAlreadyExist = true

			break
		}
	}

	if !labelAlreadyExist {
		_, _, err = c.client.Issues.CreateLabel(ctx, owner, repo, label)
		if err != nil {
			return errors.Wrap(err, "failed to create label")
		}
	}

	_, _, err = c.client.Issues.AddLabelsToIssue(ctx, owner, repo, prNum, []string{*label.Name})
	if err != nil {
		return errors.Wrap(err, "failed to add label")
	}

	return nil
}

// CreateCommentOnPullRequest created a comment on a pull request.
func (c *Client) CreateCommentOnPullRequest(ctx context.Context, owner, repo, text string, prNum int) error {
	newComment := &github.IssueComment{Body: github.String(text)}

	_, _, err := c.client.Issues.CreateComment(ctx, owner, repo, prNum, newComment)
	if err != nil {
		return errors.Wrap(err, "failed to create comment")
	}

	return nil
}

// DeleteCommentInIssue deletes a comment in issue/pull request.
func (c *Client) DeleteCommentInIssue(ctx context.Context, owner, repo string, commentID int64) error {
	_, err := c.client.Issues.DeleteComment(ctx, owner, repo, commentID)
	if err != nil {
		return errors.Wrap(err, "failed to delete comment")
	}

	return nil
}

// GetPullRequest return a pull request by number.
func (c *Client) GetPullRequest(ctx context.Context, owner, repo string, prNum int) (*github.PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(ctx, owner, repo, prNum)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get a pull request")
	}

	return pr, nil
}

// GetPullRequestReviewList returns pull request reviews by number.
func (c *Client) GetPullRequestReviewList(ctx context.Context, owner, repo string,
	prNum int) ([]*github.PullRequestReview, error) {
	list, _, err := c.client.PullRequests.ListReviews(ctx, owner, repo, prNum, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get reviews list")
	}

	return list, nil
}

// GetIssueListLabels returns labels list of an issue/pull request by number.
func (c *Client) GetIssueListLabels(ctx context.Context, owner, repo string, prNum int) ([]*github.Label, error) {
	list, _, err := c.client.Issues.ListLabelsByIssue(ctx, owner, repo, prNum, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get labels list by issue")
	}

	return list, nil
}

// CreateReview creates a review on pull request.
func (c *Client) CreateReview(ctx context.Context, owner, repo, body, event string,
	prNum int) (*github.PullRequestReview, error) {
	prReview, _, err := c.client.PullRequests.CreateReview(ctx, owner, repo, prNum,
		&github.PullRequestReviewRequest{
			Body:  &body,
			Event: &event,
		})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create review")
	}

	return prReview, nil
}

// AddAssignees adds assignees to issue/pull request.
func (c *Client) AddAssignees(ctx context.Context, owner, repo string,
	prNum int, assignees []string) (*github.Issue, error) {
	issue, _, err := c.client.Issues.AddAssignees(ctx, owner, repo, prNum, assignees)
	if err != nil {
		return nil, errors.Wrap(err, "failed to add assignees")
	}

	return issue, nil
}

// ClosePullRequest closes a pull request by number.
func (c *Client) ClosePullRequest(ctx context.Context, owner, repo string, prNum int) error {
	_, _, err := c.client.PullRequests.Edit(ctx, owner, repo, prNum, &github.PullRequest{
		State: github.String("closed"),
	})
	if err != nil {
		return errors.Wrap(err, "failed to close pull request")
	}

	return nil
}

// GetOpenPullRequestsList returns a pull request list of repository.
func (c *Client) GetOpenPullRequestsList(
	ctx context.Context, owner, repo string, perpage int,
) ([]*github.PullRequest, error) {
	pr, _, err := c.client.PullRequests.List(ctx, owner, repo, &github.PullRequestListOptions{
		ListOptions: github.ListOptions{PerPage: perpage},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get open pull requests list")
	}

	return pr, nil
}

// GetPullRequestFileList receives a pull request file list.
func (c *Client) GetPullRequestFileList(ctx context.Context, owner, repo string, prNum int,
	perpage int) ([]*github.CommitFile, error) {
	list, _, err := c.client.PullRequests.ListFiles(ctx, owner, repo, prNum, &github.ListOptions{
		PerPage: perpage,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get pull request file list")
	}

	return list, nil
}
