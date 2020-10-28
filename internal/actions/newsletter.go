package actions

import (
	e "errors"
	"log"

	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/auth"
	"github.com/txsvc/slack/pkg/slack"
)

// 1) extract callback_id from request
// 2) initiate modal, if needed
// 3) store callback_id & view ID in memory, from modal response
// 4) use view ID to lookup callback_id in order to know how to process the view submission callback

// StartAddToNewsletter starts the add to newsletter action
func StartAddToNewsletter(c *gin.Context, a *slack.ActionRequest) error {
	ctx := appengine.NewContext(c.Request)

	token, err := auth.GetToken(ctx, a.Team.ID, "slack")
	if err != nil {
		return err
	}

	// build the modal view
	m := addToNewsletterModal(a)

	var resp slack.ModalResponse
	err = slack.CustomPost(ctx, token, "views.open", &m, &resp)
	if err != nil {
		return err
	}

	if resp.OK != true {
		return errors.NewOperationError("views.open", e.New(resp.Error))
	}

	// store the correlation ID
	return slack.StoreActionCorrelation(ctx, a.CallbackID, resp.View.ID, resp.View.TeamID)
}

// CompleteAddToNewsletter completes the newsletter action
func CompleteAddToNewsletter(c *gin.Context, s *slack.ViewSubmission) error {
	// FIXME remove this
	log.Printf("s-> %v\n\n", util.PrintJSON(s))

	return nil
}

func addToNewsletterModal(a *slack.ActionRequest) *slack.ModalRequest {

	blocks := make([]interface{}, 1)

	blocks[0] = slack.SectionBlock{
		Type:    "section",
		BlockID: "block1",
		Text: slack.TextObject{
			Type: "plain_text",
			Text: a.Message.Text,
		},
	}

	m := slack.ModalRequest{
		TriggerID: a.TriggerID,
		View: slack.ViewElement{
			Type: "modal",
			Title: slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Save to newsletter",
			},
			Submit: &slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Save",
			},
			Close: &slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Cancel",
			},
			Blocks: blocks,
		},
	}

	return &m
}
