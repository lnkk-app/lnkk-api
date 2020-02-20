package actions

import (
	"log"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// StartPublishLinkAction initiates the interaction
func StartPublishLinkAction(c *gin.Context, a *slack.ActionRequest) error {
	ctx := appengine.NewContext(c.Request)

	token, err := backend.GetAuthToken(ctx, a.Team.ID)
	if err != nil {
		return err
	}

	// build the modal view
	m := newPublishLinkModal(a)
	// LOG log.Printf("modal: %v\n\n", util.PrintJSON(m))

	var resp slack.ModalResponse
	err = slack.CustomPost(c, token, "views.open", &m, &resp)
	if err != nil {
		return err
	}

	if resp.OK != true {
		return slack.NewSimpleError("views.open", resp.Error)
	}
	// FIXME remove this
	// LOG log.Printf("response: %v\n\n", util.PrintJSON(resp))

	return nil
}

// CompletePublishLinkAction completes the interaction
func CompletePublishLinkAction(c *gin.Context, s *slack.ViewSubmission) error {
	// FIXME remove this
	log.Printf("submission: %v\n\n", util.PrintJSON(s))

	return nil
}

// newPublishLinkModal creates the modal struct
func newPublishLinkModal(a *slack.ActionRequest) *slack.ModalRequest {

	blocks := make([]interface{}, 4)

	blocks[0] = slack.SectionBlock{
		Type:    "section",
		BlockID: "block1",
		Text: slack.TextObject{
			Type: "plain_text",
			Text: a.Message.Text,
		},
	}
	blocks[1] = slack.DividerBlock{
		Type: "divider",
	}
	blocks[2] = slack.InputBlock{
		Type:    "input",
		BlockID: "block2",
		Label: slack.TextObject{
			Type: "plain_text",
			Text: "Select Channels",
		},
		Element: slack.Checkboxes{
			Type:     "checkboxes",
			ActionID: "channels",
			Options: []slack.OptionsObject{
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#general",
					},
					Value: "general",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#learn",
					},
					Value: "learn",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#technology",
					},
					Value: "technology",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#events",
					},
					Value: "events",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#open_source",
					},
					Value: "open_source",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "#random",
					},
					Value: "random",
				},
			},
		},
	}

	blocks[3] = slack.InputBlock{
		Type:    "input",
		BlockID: "block3",
		Label: slack.TextObject{
			Type: "plain_text",
			Text: "Publish ...",
		},
		Element: slack.Radiobuttons{
			Type:     "radio_buttons",
			ActionID: "buttons",
			Options: []slack.OptionsObject{
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "now",
					},
					Value: "now",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "today",
					},
					Value: "today",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "tomorrow",
					},
					Value: "tomorrow",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "this week",
					},
					Value: "this_week",
				},
				{
					Text: slack.TextObject{
						Type: "plain_text",
						Text: "next week",
					},
					Value: "next_week",
				},
			},
			InitialOption: &slack.OptionsObject{
				Text: slack.TextObject{
					Type: "plain_text",
					Text: "now",
				},
				Value: "now"},
		},
	}

	m := slack.ModalRequest{
		TriggerID: a.TriggerID,
		View: slack.ViewElement{
			Type: "modal",
			Title: slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Publish",
			},
			Submit: &slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Publish",
			},
			Close: &slack.DefaultViewElement{
				Type: "plain_text",
				Text: "Close",
			},
			Blocks: blocks,
		},
	}

	return &m
}
