package cmd

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/slack/pkg/slack"
)

// SlackCmdShorten handles the /lnkk group of commands
func SlackCmdShorten(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("Command: '%s' %s", cmd.Command, cmd.Txt),
				},
			},
		},
	}, nil

}
