package api

import (
	e "errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"

	"google.golang.org/appengine"
)

// OAuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token. See https://api.slack.com/docs/oauth
func OAuthEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// extract parameters
	code := c.Query("code")
	redirectURI := c.Query("redirect_uri")

	// FIXME secure the request by using a state
	// state := c.Query("state")

	if code != "" {
		// exchange the temporary code with a real auth token
		resp, err := slack.GetOAuthToken(ctx, code)

		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		// get team info
		var teamInfo slack.TeamInfo
		err = slack.Get(ctx, resp.AccessToken, "team.info", "", &teamInfo)
		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		if teamInfo.OK == false {
			platform.Report(errors.NewOperationError("team.info", e.New(teamInfo.Error)))
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		err = backend.UpdateAuthorization(ctx, teamInfo.Team.ID, teamInfo.Team.Name, resp.AccessToken, resp.TokenType, resp.Scope, resp.AppID, resp.BotUserID)
		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}
	}

	// back to the sign-up process on the main website
	if redirectURI == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	}
}
