package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/modal"
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
		err = publishLinkAction(c, &action)
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

		// FIXME remove this
		log.Printf("submission: %v\n\n", util.PrintJSON(submission))
	} else {
		// FIXME should not happen !
	}
}

func publishLinkAction(c *gin.Context, a *slack.ActionRequest) error {
	ctx := appengine.NewContext(c.Request)

	token, err := backend.GetAuthToken(ctx, a.Team.ID)
	if err != nil {
		return err
	}

	// build the modal view
	m := modal.CreatePublishLinkModal(a)
	// LOG log.Printf("modal: %v\n\n", util.PrintJSON(m))

	var resp slack.ModalResponse
	err = slack.CustomPost(c, token, "views.open", &m, &resp)
	if err != nil {
		return err
	}

	// FIXME remove this
	// LOG log.Printf("response: %v\n\n", util.PrintJSON(resp))

	return nil
}
