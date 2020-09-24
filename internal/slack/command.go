package slack

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/txsvc/commons/pkg/env"
	"google.golang.org/appengine"

	"github.com/lnkk-app/lnkk-api/pkg/shortener"
	"github.com/lnkk-app/lnkk-api/pkg/slack"
	s "github.com/lnkk-app/lnkk-api/pkg/slack"
	"github.com/lnkk-app/lnkk-api/pkg/types"
)

// CmdLnkkHandler dispatches /lnkk commands
func CmdLnkkHandler(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	if len(cmd.Txt) == 0 {
		return shortHelpMessage(), nil
	}
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) == 0 {
		return shortHelpMessage(), nil
	}

	subcmd := parts[0]
	if subcmd == "help" {
		return helpMessage(), nil
	} else if subcmd == "shorten" {
		return handleShorten(c, cmd)
	} else if subcmd == "stats" {
		return handleStats(c, cmd)
	} else if subcmd == "subscribe" {
		return handleSubscribe(c, cmd)
	} else if subcmd == "feed" {
		return handleFeed(c, cmd)
	}
	return shortHelpMessage(), nil
}

func handleShorten(c *gin.Context, cmd *s.SlashCommand) (*s.SectionBlocks, error) {
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) < 2 {
		return &s.SectionBlocks{
			Blocks: []s.SectionBlock{
				{
					Type: "section",
					Text: s.TextObject{
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

	return &s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: msg,
				},
			},
		},
	}, nil
}

func handleStats(c *gin.Context, cmd *s.SlashCommand) (*s.SectionBlocks, error) {
	return &s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk stats",
				},
			},
		},
	}, nil
}

func handleSubscribe(c *gin.Context, cmd *s.SlashCommand) (*s.SectionBlocks, error) {
	return &s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk subscribe",
				},
			},
		},
	}, nil
}

func handleFeed(c *gin.Context, cmd *s.SlashCommand) (*s.SectionBlocks, error) {
	return &s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed",
				},
			},
		},
	}, nil
}

func helpMessage() *s.SectionBlocks {
	response := s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk shorten URL\nRetunrs a shortened version of the url",
				},
			},
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk stats [ID]\nReturns stats for a url. If ID is omitted, returns stats for all urls in the past 24h",
				},
			},
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk subscribe URL\nSubscribes to a RSS feed",
				},
			},
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed list\nLists all RSS subscriptions for this channel",
				},
			},
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk feed remove ID\nRemoves a RSS subscription from this channel",
				},
			},
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "/lnkk help\nThis command",
				},
			},
		},
	}

	return &response
}

func shortHelpMessage() *s.SectionBlocks {
	return &s.SectionBlocks{
		Blocks: []s.SectionBlock{
			{
				Type: "section",
				Text: s.TextObject{
					Type: "mrkdwn",
					Text: "Unsupported command.\nUse */lnkk help* to learn more about all supported commands",
				},
			},
		},
	}
}
