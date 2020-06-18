package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/internal/command"
	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// SlashCmdEndpoint receives callbacks from Slack command /lnkk
func SlashCmdEndpoint(c *gin.Context) {

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
							Text: "Oops, something went wrong",
						},
					},
				},
			}
			resp = &response
		}
	}

	c.JSON(status, resp)
}

// FIXME make this a handler function similar to the ActionRequests
func handleCmdLnkk(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	if len(cmd.Txt) == 0 {
		return command.ShortHelpMessage(), nil
	}
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) == 0 {
		return command.ShortHelpMessage(), nil
	}

	subcmd := parts[0]
	if subcmd == "help" {
		return command.HelpMessage(), nil
	} else if subcmd == "shorten" {
		return command.HandleSubcmdShorten(c, cmd)
	} else if subcmd == "stats" {
		return command.HandleSubcmdStats(c, cmd)
	} else if subcmd == "subscribe" {
		return command.HandleSubcmdSubscribe(c, cmd)
	} else if subcmd == "feed" {
		return command.HandleSubcmdFeed(c, cmd)
	}
	return command.ShortHelpMessage(), nil
}
