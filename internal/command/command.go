package command

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/env"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/pkg/shortener"
	"github.com/lnkk-ai/lnkk/pkg/slack"
	"github.com/lnkk-ai/lnkk/pkg/types"
)

func HandleSubcmdShorten(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) < 2 {
		return &slack.SectionBlocks{
			Blocks: []slack.SectionBlock{
				{
					Type: "section",
					Text: slack.TextObject{
						Type: "mrkdwn",
						Text: "Usage: /lnkk shorten *URL*",
					},
				},
			},
		}, nil
	}

	msg := ""
	tags := ""
	if len(parts) > 2 {
		tags = strings.Join(parts[2:], ",")
	}
	asset := types.Asset{
		URL:       parts[1],
		Owner:     cmd.UserID,
		Source:    "slack",
		Client:    "lnkk",
		Affiliate: cmd.TeamID,
		Tags:      tags,
	}

	ctx := appengine.NewContext(c.Request)
	uri, err := shortener.CreateAsset(ctx, &asset)

	if err != nil {
		msg = err.Error() // FIXME better error message
	} else {
		url := env.Getenv("SHORT_URL", "https://lnkk.host")
		msg = fmt.Sprintf("%s/r/%s", url, uri)
	}

	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: msg,
				},
			},
		},
	}, nil
}

func HandleSubcmdStats(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk stats",
				},
			},
		},
	}, nil
}

func HandleSubcmdSubscribe(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk subscribe",
				},
			},
		},
	}, nil
}

func HandleSubcmdFeed(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed",
				},
			},
		},
	}, nil
}

func HelpMessage() *slack.SectionBlocks {
	response := slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk shorten URL\nRetunrs a shortened version of the url",
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

func ShortHelpMessage() *slack.SectionBlocks {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "Unsupported command.\nUse */lnkk help* to learn more about all supported commands",
				},
			},
		},
	}
}
