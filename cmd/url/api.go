package main

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"

	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// redirectEndpoint receives a URI to be shortened
func redirectEndpoint(c *gin.Context) {
	//topic := "api.redirect.get"
	ctx := appengine.NewContext(c.Request)

	uri := c.Param("uri")
	if uri == "" {
		// TODO log this event
		// FIXME redirect to a valid page
		c.String(http.StatusOK, "42")
		return
	}

	a, err := backend.GetAsset(ctx, uri)
	if err != nil {
		// TODO log this event
		// FIXME redirect to a valid page
		c.String(http.StatusOK, "42")
		return
	}

	// audit, i.e. extract some user data
	now := util.Timestamp()
	m := types.MeasurementDS{
		URI:            uri,
		User:           "anonymous",
		IP:             c.ClientIP(),
		UserAgent:      strings.ToLower(c.GetHeader("User-Agent")),
		AcceptLanguage: strings.ToLower(c.GetHeader("Accept-Language")),
		Day:            util.TimestampToWeekday(now),
		Hour:           util.TimestampToHour(now),
		Created:        now,
	}
	backend.CreateMeasurement(ctx, &m)

	// TODO log the event
	c.Redirect(http.StatusTemporaryRedirect, a.URL)
}

// slashCmdShortenEndpoint receives a URL to be shortened
func slashCmdShortenEndpoint(c *gin.Context) {
	//topic := "slack.shorten.post"
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
}

/*

// AsExternal create an external representation of the asset
func (t *AssetDS) AsExternal() *api.Asset {
	asset := api.Asset{
		URI:       t.URI,
		URL:       t.URL,
		SecretID:  t.SecretID,
		Cohort:    t.Cohort,
		Affiliate: t.Affiliate,
		Tags:      t.Tags,
	}
	return &asset
}

*/
