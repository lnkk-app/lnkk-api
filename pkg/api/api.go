package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
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
		// TODO log this event
		c.String(http.StatusOK, "42")
		return
	}

	a, err := urlshortener.GetURL(ctx, shortLink)
	if err != nil {
		// TODO log this event
		c.String(http.StatusOK, "42") // FIXME this is stupid ...
		return
	}

	// log the event and redirect
	urlshortener.LogRedirectRequest(ctx, shortLink, c)
	c.Redirect(http.StatusTemporaryRedirect, a.Link)
}
