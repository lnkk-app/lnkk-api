package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/pkg/slack"
)

/*

ctx := appengine.NewContext(c.Request)

	params := slack.SlashCommandPayload(c)

	uri, _ := util.ShortUUID()
	secret, _ := util.ShortUUID()
	source, _ := params["team_id"]
	owner, _ := params["user_id"]

	// parse the text into the url and tags if there are any
	// example: https://foo.bar.com foo bar
	text, _ := params["text"]
	parts := strings.Split(text, " ")
	url := parts[0]
	tags := ""
	if len(parts) > 1 {
		tags = strings.Join(parts[1:], ",")
	}

	asset := types.AssetDS{
		URI:       uri,
		URL:       url,
		Owner:     owner,
		SecretID:  secret,
		Source:    source,
		Cohort:    "slack",
		Affiliate: "",
		Tags:      tags,
	}

	err := backend.CreateAsset(ctx, &asset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}

	shortened := fmt.Sprintf("%s/r/%s", env.Getenv("BASE_URL", "/"), uri)
	c.JSON(http.StatusOK, shortened)

*/

// CmdLnkkEndpoint receives callbacks from Slack command /lnkk
func CmdLnkkEndpoint(c *gin.Context) {

	//ctx := appengine.NewContext(c.Request)

	cmd := slack.GetSlashCommand(c)

	//uri, _ := util.ShortUUID()
	//secret, _ := util.ShortUUID()

	s := fmt.Sprintf("%s,%s %s/%s %s - %s", cmd.TeamID, cmd.UserName, cmd.ChannelID, cmd.ChannelName, cmd.Command, cmd.Txt)
	response := slack.SectionBlocks{
		Blocks: make([]slack.SectionBlock, 1),
	}
	response.Blocks[0] = slack.SectionBlock{
		Type: "section",
		Text: slack.TextObject{
			Type: "mrkdwn",
			Text: s,
		},
	}

	c.JSON(http.StatusOK, &response)
}
