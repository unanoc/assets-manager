package http

import (
	"errors"
	"net/http"

	ghlib "github.com/google/go-github/v38/github"
	log "github.com/sirupsen/logrus"

	"github.com/trustwallet/assets-manager/internal/config"
	"github.com/trustwallet/assets-manager/internal/services/worker/events"
)

const (
	eventActionCreated     = "created"
	eventActionOpened      = "opened"
	eventActionReopened    = "reopened"
	eventActionSynchronize = "synchronize"
)

var ErrUnknownEvent = errors.New("unknown event")

// GithubEventsHandler returns http handler for handling github events.
func GithubEventsHandler(eh *events.EventHandler) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Validate webhook payload.
		payloadBytes, err := ghlib.ValidatePayload(r, []byte(config.Default.Github.AppWebhookSecret))
		if err != nil {
			sendHTTPErrorResponse(w, err, http.StatusForbidden)

			return
		}

		// Parse the incoming request into an event.
		eventPayload, err := ghlib.ParseWebHook(ghlib.WebHookType(r), payloadBytes)
		if err != nil {
			sendHTTPErrorResponse(w, err, http.StatusBadRequest)

			return
		}

		log.WithField("event", ghlib.WebHookType(r)).Debug("Incoming event")

		// Handle the incoming event.
		switch event := eventPayload.(type) {
		case *ghlib.PullRequestEvent:
			if *event.Action == eventActionOpened || *event.Action == eventActionReopened {
				err = eh.HandlePullRequestOpened(r.Context(), event)
			}

			if *event.Action == eventActionSynchronize {
				err = eh.HandlePullRequestChangesPushed(r.Context(), event)
			}

		case *ghlib.IssueCommentEvent:
			if *event.Action == eventActionCreated {
				err = eh.HandleIssueCommentCreated(r.Context(), event)
			}

		case *ghlib.PullRequestReviewCommentEvent:
			if *event.Action == eventActionOpened {
				err = eh.HandlePullRequestReviewCommentCreated(r.Context(), event)
			}

		default:
			sendHTTPErrorResponse(w, ErrUnknownEvent, http.StatusBadRequest)

			return
		}

		if err != nil {
			sendHTTPErrorResponse(w, err, http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func sendHTTPErrorResponse(w http.ResponseWriter, err error, status int) {
	if err != nil {
		log.Error(err)
	}

	http.Error(w, err.Error(), status)
}
