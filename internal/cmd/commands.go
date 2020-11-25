package cmd

import (
	"context"
	e "errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/slack/pkg/slack"
)

const (
	helpText = "help [or shorten, stats, top, delete ...]\n\n" +
		"*shorten* _long url_: shortens a long URL\n" +
		"*list* _n_: returns the last _n_ assets\n" +
		"*stats* [_token_]: returns last 24h stats or all time stats for a specific url\n" +
		"*top*: all time top 10 of urls\n" +
		"*delete* _token_: deletes the specified url.\n\n" +
		"For more information see https://lnkk.host/help"
)

func init() {
	RegisterCmdHandler("lnkk", "shorten", shortenCmdHandler)
	RegisterCmdHandler("lnkk", "list", listCmdHandler)
	RegisterCmdHandler("lnkk", "stats", statsCmdHandler)
	RegisterCmdHandler("lnkk", "top", topCmdHandler)
	RegisterCmdHandler("lnkk", "delete", deleteCmdHandler)
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
		Link:   cmdLn[0],
		Owner:  ownerID(cmd),
		Source: strings.ToLower("slack." + cmd.TeamID),
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

func listCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	var txt string

	count := 10
	if len(cmdLn) == 1 {
		i, err := strconv.ParseInt(cmdLn[0], 10, 64)
		if err != nil {
			platform.ReportError(err)
			return somethingWentWrong(), nil
		}
		count = (int)(i)
	}
	if count < 1 {
		count = 10
	}

	iter := urlshortener.GetAssets(ctx, ownerID(cmd), count, 0)
	for {
		var asset urlshortener.Asset
		_, err := iter.Next(&asset)
		if err == iterator.Done {
			break
		}
		if err != nil {
			platform.ReportError(err)
			return somethingWentWrong(), nil
		}
		txt = txt + fmt.Sprintf("%s\n%s Token: %v\n\n", asset.LongLink, asset.PreviewLink, asset.AccessToken)
	}

	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: txt,
				},
			},
		},
	}, nil
}

func statsCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	since := util.Timestamp() - int64(24*60*60)

	assets, err := urlshortener.NewAssetsSince(ctx, ownerID(cmd), since)
	if err != nil {
		platform.ReportError(err)
		assets = 0
	}

	redirects, err := urlshortener.RedirectsSince(ctx, ownerID(cmd), since)
	if err != nil {
		platform.ReportError(err)
		redirects = 0
	}

	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: fmt.Sprintf("Statistics for the last 24 hours:\n\nAssets created: %v\nRedirects: %v", assets, redirects),
				},
			},
		},
	}, nil
}

func topCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	return notImplemented(), nil
}

func deleteCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	return notImplemented(), nil
}

func helpCmdHandler(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error) {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: helpText,
				},
			},
		},
	}, nil
}

func ownerID(c *slack.SlashCommand) string {
	return strings.ToLower(c.TeamID + "." + c.UserID)
}

func notImplemented() *slack.SectionBlocks {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "Sorry, not implemented yet.",
				},
			},
		},
	}
}

func somethingWentWrong() *slack.SectionBlocks {
	return &slack.SectionBlocks{
		Blocks: []slack.SectionBlock{
			{
				Type: "section",
				Text: slack.TextObject{
					Type: "mrkdwn",
					Text: "Oops, something went wrong.",
				},
			},
		},
	}
}
