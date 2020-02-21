package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"

	"github.com/lnkk-ai/lnkk/internal/actions"
)

// ActionRequestEndpoint receives callbacks from Slack
func ActionRequestEndpoint(c *gin.Context) {
	var peek slack.AppRequestPeek

	err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &peek)
	if err != nil {
		// report the error and abort
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

		//LOG log.Printf("action: %v\n\n", util.PrintJSON(action))

		// FIXME switch actions
		err = actions.StartPublishLinkAction(c, &action)
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

		actions.CompletePublishLinkAction(c, &submission)
	} else {
		platform.Report(errors.New(fmt.Sprint("Unknown action request: '%s", peek.Type)))
	}
}
