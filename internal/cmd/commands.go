package cmd

import (
	"context"
	e "errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
	"google.golang.org/appengine"

	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/slack/pkg/slack"
)

func init() {
	RegisterCmdHandler("lnkk", "shorten", shortenCmdHandler)
	RegisterCmdHandler("lnkk", "help", helpCmdHandler)
}

// SlackCmdLnkk handles the /lnkk group of commands
func SlackCmdLnkk(c *gin.Context, cmd *slack.SlashCommand) (*slack.SectionBlocks, error) {
	ctx := appengine.NewContext(c.Request)

	resp, err := ParseCommandLine(ctx, "lnkk", cmd)
	if err != nil {
		platform.ReportError(err)
		return nil, slack.NewSlackCmdEror(err.Error(), cmd, err)
	}
	return resp.(*slack.SectionBlocks), nil
}

func shortenCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	if len(cmdLn) == 0 || cmdLn == nil {
		return nil, e.New("Missing URL")
	}

	asset := urlshortener.AssetRequest{
		Link:  cmdLn[0],
		Owner: strings.ToLower(cmd.TeamID + "." + cmd.UserID),
	}
	_asset, err := urlshortener.CreateURL(ctx, &asset)
	if err != nil {
		return nil, err
	}

	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("Link: %s\nShort: %s\nToken: %s", _asset.Link, _asset.PreviewLink, _asset.AccessToken),
				},
			},
		},
	}, nil
}

func helpCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	return nil, nil
}
