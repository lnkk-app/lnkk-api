package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// CmdLnkkEndpoint receives callbacks from Slack command /lnkk
func CmdLnkkEndpoint(c *gin.Context) {

	status := http.StatusOK

	// extract the cmd and react on it
	cmd := slack.GetSlashCommand(c)
	resp, err := handleCmdLnkk(c, cmd)

	if err != nil {
		status = http.StatusBadRequest
		if resp == nil {
			response := slack.SectionBlocks{
				Blocks: []slack.SectionBlock{
					{
						Type: "section",
						Text: slack.TextObject{
							Type: "mrkdwn",
							Text: "Oops, something happened",
						},
					},
				},
			}
			resp = &response
		}

		c.JSON(status, resp)
		return
	}

	if resp == nil {
		response := slack.SectionBlocks{
			Blocks: []slack.SectionBlock{
				{
					Type: "section",
					Text: slack.TextObject{
						Type: "mrkdwn",
						Text: "This should not happen",
					},
				},
			},
		}

		resp = &response
		status = http.StatusBadRequest
	}

	c.JSON(status, resp)

}

func handleCmdLnkk(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	if len(cmd.Txt) == 0 {
		return helpShortMessage(), nil
	}
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) == 0 {
		return helpShortMessage(), nil
	}

	subcmd := parts[0]
	if subcmd == "help" {
		return helpMessage(), nil
	} else if subcmd == "shorten" {

	}

	return helpShortMessage(), nil
}

func helpMessage() *slack.SectionBlocks {
	response := slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk shorten URL\nShortens a url and",
				},
			},
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk stats [ID]\nReturns stats for a url. If ID is omitted, returns stats for all urls in the past 24h",
				},
			},
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk subscribe URL\nSubscribes to a RSS feed",
				},
			},
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed list\nLists all RSS subscriptions for this channel",
				},
			},
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed remove ID\nRemoves a RSS subscription from this channel",
				},
			},
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk help\nThis command",
				},
			},
		},
	}

	return &response
}

func helpShortMessage() *slack.SectionBlocks {
	response := slack.SectionBlocks{
		Blocks: make([]slack.SectionBlock, 1),
	}
	response.Blocks[0] = slack.SectionBlock{
		Type: "section",
		Text: slack.TextObject{
			Type: "mrkdwn",
			Text: "Unsupported command.\nUse */lnkk help* to learn more about all supported commands",
		},
	}
	return &response
}
