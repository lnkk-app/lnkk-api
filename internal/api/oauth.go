package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/lnkk-ai/lnkk/pkg/logger"
	"github.com/lnkk-ai/lnkk/pkg/slack"
	"github.com/lnkk-ai/lnkk/pkg/tasks"
	"github.com/majordomusio/commons/pkg/env"
	"google.golang.org/appengine"
)

// AuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token.
// See https://api.slack.com/docs/oauth
func AuthEndpoint(c *gin.Context) {
	topic := "api.oauth"
	ctx := appengine.NewContext(c.Request)

	// extract parameters
	code := c.Query("code")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")

	logger.Info(topic, "code=%s, state=%s", code, state)

	if code == "" {
		// OAuth was cancelled, shoud check for error=access_denied
		logger.Warn(topic, c.Query("error"))
	} else {
		// exchange the temporary code with a real auth token
		resp, err := slack.OAuthAccess(ctx, code)

		if err != nil {
			errorreporting.Report(err)
		} else {

			// FIXME error handling ?
			backend.UpdateAuthorization(ctx, resp.TeamID, resp.TeamName, resp.AccessToken, resp.Scope, resp.AuthorizingUser.UserID, resp.InstallerUser.UserID)
			backend.UpdateWorkspace(ctx, resp.TeamID, resp.TeamName)

			// schedule the first update of the new workspace
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/users?id="+resp.TeamID)
			tasks.Schedule(ctx, backend.BackgroundWorkQueue, env.Getenv("BASE_URL", "")+api.JobsBaseURL+"/channels?id="+resp.TeamID)

			backend.MarkWorkspaceUpdated(ctx, resp.TeamID)
			logger.Info(topic, "workspace=%s", resp.TeamID)
		}
	}

	// back to the sign-up process on the main website
	if redirectURI == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	}
}
