package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/pkg/shortener"
)

const (
	// APIPrefix is the common namespace prefix of API releated routes
	APIPrefix string = "/a/1"
	// CronBaseURL is the namespace for all scheduler routes
	CronBaseURL = "/_c/1"
	// WorkerBaseURL is the namespace for all worker routes
	WorkerBaseURL = "/_w/1"
)

// ShortenEndpoint receives a URI to be shortened
func ShortenEndpoint(c *gin.Context) {
	//topic := "api.shorten.post"
	ctx := appengine.NewContext(c.Request)

	var asset shortener.AssetRequest
	err := c.BindJSON(&asset)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	_asset, err := shortener.CreateURL(ctx, &asset)
	svc.StandardJSONResponse(c, _asset, err)
}

// RedirectEndpoint receives a URI to be shortened
func RedirectEndpoint(c *gin.Context) {
	//topic := "api.redirect.get"
	ctx := appengine.NewContext(c.Request)

	shortLink := c.Param("short")
	if shortLink == "" {
		// FIXME: log this event
		c.Redirect(http.StatusTemporaryRedirect, env.GetString("BASE_URL", "https://lnkk.host"))
		return
	}

	asset, err := shortener.GetURL(ctx, shortLink, true)
	if err != nil || asset.State != shortener.StateActive {
		// FIXME: log this event
		redirectToErrorPage := fmt.Sprintf("%s/e/%s", env.GetString("BASE_URL", "https://lnkk.host"), shortLink)
		c.Redirect(http.StatusTemporaryRedirect, redirectToErrorPage)
		return
	}

	// log the event and redirect
	shortener.LogRedirectRequest(ctx, asset, c)
	c.Redirect(http.StatusTemporaryRedirect, asset.LongLink)
}
