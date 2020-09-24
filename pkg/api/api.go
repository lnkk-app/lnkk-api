package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/pkg/shortener"
	"github.com/lnkk-app/lnkk-api/pkg/types"
)

// ShortenEndpoint receives a URI to be shortened
func ShortenEndpoint(c *gin.Context) {
	//topic := "api.shorten.post"
	ctx := appengine.NewContext(c.Request)

	var asset types.Asset
	var uri string

	err := c.BindJSON(&asset)
	if err == nil {
		uri, err = shortener.CreateAsset(ctx, &asset)
		asset.URI = uri
	}

	svc.StandardJSONResponse(c, asset, err)
}

// RedirectEndpoint receives a URI to be shortened
func RedirectEndpoint(c *gin.Context) {
	//topic := "api.redirect.get"
	ctx := appengine.NewContext(c.Request)

	uri := c.Param("uri")
	if uri == "" {
		// TODO log this event
		c.String(http.StatusOK, "42")
		return
	}

	a, err := shortener.GetAsset(ctx, uri)
	if err != nil {
		// TODO log this event
		c.String(http.StatusOK, "42")
		return
	}

	// audit, i.e. extract some user data
	m := shortener.MeasurementDS{ // FIXME use a public struct
		URI:            uri,
		User:           "anonymous",
		IP:             c.ClientIP(),
		UserAgent:      strings.ToLower(c.GetHeader("User-Agent")),
		AcceptLanguage: strings.ToLower(c.GetHeader("Accept-Language")),
		Created:        util.Timestamp(),
	}
	shortener.CreateMeasurement(ctx, &m)

	// TODO log the event
	c.Redirect(http.StatusTemporaryRedirect, a.URL)
}
