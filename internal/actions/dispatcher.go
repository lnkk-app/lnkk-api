package actions

import (
	"context"
	e "errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"
	"google.golang.org/appengine"
)

//  map[string]string

// StartActionFunc is a callback for starting a action
type StartActionFunc func(*gin.Context, *slack.ActionRequest) error

// CompleteActionFunc is a callback for completing an action
type CompleteActionFunc func(*gin.Context, *slack.ViewSubmission) error

var startActionLookup map[string]StartActionFunc
var completeActionLookup map[string]CompleteActionFunc

func init() {
	// initialize the start action lookup table
	startActionLookup = make(map[string]StartActionFunc, 1)
	startActionLookup["publish_lnkk"] = StartPublishLinkAction

	// initialize the complete action lookup table
	completeActionLookup = make(map[string]CompleteActionFunc, 1)
	completeActionLookup["publish_lnkk"] = CompletePublishLinkAction
}

// StartAction initiates a dialog with the user
func StartAction(c *gin.Context, a *slack.ActionRequest) error {
	// 1) extract callback_id from request
	// 2) initiate modal, if needed
	// 3) store callback_id & view ID in memory, from modal response

	action := a.CallbackID
	handler := startActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action request '%s'", action)))
	}

	return handler(c, a)
}

// CompleteAction starts the processing of the action's result
func CompleteAction(c *gin.Context, s *slack.ViewSubmission) error {
	// 4) use view ID to lookup callback_id in order to know how to process the view submission callback

	ctx := appengine.NewContext(c.Request)
	action := lookupActionCorrelation(ctx, s.View.ID, s.Team.ID)

	if action == "" {
		return nil
	}

	handler := completeActionLookup[action]
	if handler == nil {
		return errors.NewOperationError(action, e.New(fmt.Sprintf("No handler for action response '%s'", action)))
	}

	return handler(c, s)
}

func storeActionCorrelation(ctx context.Context, action, viewID, teamID string) error {
	key := viewID + teamID
	return platform.Set(ctx, key, action, "15m")
}

func lookupActionCorrelation(ctx context.Context, viewID, teamID string) string {
	key := viewID + teamID
	v, _ := platform.Get(ctx, key)
	return v
}
