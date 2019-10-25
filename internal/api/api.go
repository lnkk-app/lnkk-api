package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-ai/lnkk/internal/store"
	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/majordomusio/commons/pkg/util"
	"google.golang.org/appengine"
)

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

	a, err := store.GetAsset(ctx, uri)
	if err != nil {
		// TODO log this event
		c.String(http.StatusOK, "42")
		return
	}

	// audit, i.e. extract some user data
	m := types.MeasurementDS{
		URI:            uri,
		User:           "anonymous",
		IP:             c.ClientIP(),
		UserAgent:      strings.ToLower(c.GetHeader("User-Agent")),
		AcceptLanguage: strings.ToLower(c.GetHeader("Accept-Language")),
		Created:        util.Timestamp(),
	}
	store.CreateMeasurement(ctx, &m)

	// TODO log the event
	c.Redirect(http.StatusTemporaryRedirect, a.URL)
}
