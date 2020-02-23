package api

import (
	"encoding/json"
	e "errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"

	"github.com/lnkk-ai/lnkk/internal/actions"
)

// StartActionFunc is a callback for starting a action
type StartActionFunc func(*gin.Context, *slack.ActionRequest) error

// CompleteActionFunc is a callback for completing an action
type CompleteActionFunc func(*gin.Context, *slack.ViewSubmission) error

var startActionLookup map[string]StartActionFunc
var completeActionLookup map[string]CompleteActionFunc

func init() {
	// initialize the start action lookup table
	startActionLookup = make(map[string]StartActionFunc, 1)
	startActionLookup["add_newsletter"] = actions.StartAddToNewsletter

	// initialize the complete action lookup table
	completeActionLookup = make(map[string]CompleteActionFunc, 1)
	completeActionLookup["add_newsletter"] = actions.CompleteAddToNewsletter
}

// ActionRequestEndpoint receives callbacks from Slack
func ActionRequestEndpoint(c *gin.Context) {
	var peek slack.ActionRequestPeek

	err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &peek)
	if err != nil {
		platform.Report(err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
		return
	}

	if peek.Type == "message_action" {
		var action slack.ActionRequest
		err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &action)
		if err != nil {
			platform.Report(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

		err = startAction(c, &action)
		if err != nil {
			platform.Report(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

	} else if peek.Type == "view_submission" {
		var submission slack.ViewSubmission
		err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &submission)
		if err != nil {
			platform.Report(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}

		err = completeAction(c, &submission)
		if err != nil {
			platform.Report(err)
			c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "msg": err.Error()})
			return
		}
	} else {
		platform.Report(fmt.Errorf("Unknown action request: '%s'", peek.Type))
	}
}

// startAction initiates a dialog with the user
func startAction(c *gin.Context, a *slack.ActionRequest) error {
	action := a.CallbackID
	handler := startActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action request '%s'", action)))
	}

	return handler(c, a)
}

// completeAction starts the processing of the action's result
func completeAction(c *gin.Context, s *slack.ViewSubmission) error {
	ctx := appengine.NewContext(c.Request)

	action := actions.LookupActionCorrelation(ctx, s.View.ID, s.Team.ID)
	log.Printf("action -> %s\n\n", action)
	if action == "" {
		return nil
	}

	handler := completeActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action response '%s'", action)))
	}

	return handler(c, s)
}
