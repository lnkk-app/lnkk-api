package api

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// AuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token.
// See https://api.slack.com/docs/oauth
func AuthEndpoint(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	// extract parameters
	code := c.Query("code")
	redirectURI := c.Query("redirect_uri")
	// FIXME state := c.Query("state")

	if code != "" {
		// exchange the temporary code with a real auth token
		resp, err := slack.OAuthAccess(ctx, code)

		if err != nil {
			errorreporting.Report(err)
		} else {

			// FIXME error handling ?
			backend.UpdateAuthorization(ctx, resp.TeamID, resp.TeamName, resp.AccessToken, resp.Scope, resp.AuthorizingUser.UserID, resp.InstallerUser.UserID)
			backend.UpdateWorkspace(ctx, resp.TeamID, resp.TeamName)

			// schedule the first update of the new workspace
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+JobsBaseURL+"/users?id="+resp.TeamID)
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+JobsBaseURL+"/channels?id="+resp.TeamID)

			backend.MarkWorkspaceUpdated(ctx, resp.TeamID)
		}
	}

	// back to the sign-up process on the main website
	if redirectURI == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	}
}
