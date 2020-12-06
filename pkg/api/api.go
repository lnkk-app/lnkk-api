package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
)

const (
	// APIPrefix is the common namespace prefix of API releated routes
	APIPrefix string = "/a/1"
	// SchedulerBaseURL is the namespace for all scheduler routes
	SchedulerBaseURL = "/_i/1"
	// WorkerBaseURL is the namespace for all worker routes
	WorkerBaseURL = "/_w/1"
)

// ShortenEndpoint receives a URI to be shortened
func ShortenEndpoint(c *gin.Context) {
	//topic := "api.shorten.post"
	ctx := appengine.NewContext(c.Request)

	var asset urlshortener.AssetRequest
	err := c.BindJSON(&asset)
	if err != nil {
		svc.StandardJSONResponse(c, nil, err)
		return
	}

	_asset, err := urlshortener.CreateURL(ctx, &asset)
	svc.StandardJSONResponse(c, _asset, err)
}

// RedirectEndpoint receives a URI to be shortened
func RedirectEndpoint(c *gin.Context) {
	//topic := "api.redirect.get"
	ctx := appengine.NewContext(c.Request)

	shortLink := c.Param("short")
	if shortLink == "" {
		// FIXME: log this event
		c.Redirect(http.StatusTemporaryRedirect, env.Getenv("BASE_URL", "https://lnkk.host"))
		return
	}

	asset, err := urlshortener.GetURL(ctx, shortLink, true)
	if err != nil || asset.State != urlshortener.StateActive {
		// FIXME: log this event
		redirectToErrorPage := fmt.Sprintf("%s/e/%s", env.Getenv("BASE_URL", "https://lnkk.host"), shortLink)
		c.Redirect(http.StatusTemporaryRedirect, redirectToErrorPage)
		return
	}

	// log the event and redirect
	urlshortener.LogRedirectRequest(ctx, asset, c)
	c.Redirect(http.StatusTemporaryRedirect, asset.LongLink)
}
