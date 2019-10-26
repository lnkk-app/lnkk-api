package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"
)

// CmdShortenEndpoint receives a URL to be shortened
func CmdShortenEndpoint(c *gin.Context) {
	//topic := "slack.shorten.post"
	ctx := appengine.NewContext(c.Request)

	params := slashCommandPayload(c)

	uri, _ := util.ShortUUID()
	secret, _ := util.ShortUUID()
	source, _ := params["team_id"]
	owner, _ := params["user_id"]
	url, _ := params["text"]

	asset := types.AssetDS{
		URI:       uri,
		URL:       url,
		Owner:     owner,
		SecretID:  secret,
		Source:    source,
		Cohort:    "slack",
		Affiliate: "",
		Tags:      "",
	}

	err := backend.CreateAsset(ctx, &asset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}

	shortened := fmt.Sprintf("https://lnkk.host/r/%s", uri)
	c.JSON(http.StatusOK, shortened)
}

func slashCommandPayload(c *gin.Context) map[string]string {
	p := make(map[string]string)

	p["team_id"] = c.PostForm("team_id")
	p["team_domain"] = c.PostForm("team_domain")
	p["enterprise_id"] = c.PostForm("enterprise_id")
	p["enterprise_name"] = c.PostForm("enterprise_name")
	p["channel_id"] = c.PostForm("channel_id")
	p["channel_name"] = c.PostForm("channel_name")
	p["user_id"] = c.PostForm("user_id")
	p["user_name"] = c.PostForm("user_name")
	p["command"] = c.PostForm("command")
	p["text"] = c.PostForm("text")
	p["response_url"] = c.PostForm("response_url")
	p["trigger_id"] = c.PostForm("trigger_id")

	return p
}
