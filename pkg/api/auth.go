package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"

	"github.com/majordomusio/commons/pkg/util"

	"google.golang.org/appengine"
)

// AuthEndpoint handles the callback from Slack with the temporary access code
// and exchanges it with the real auth token. See https://api.slack.com/docs/oauth
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
			platform.Report(err)
		} else {
			// FIXME error handling ?
			updateAuthorization(ctx, resp.TeamID, resp.TeamName, resp.AccessToken, resp.Scope, resp.AuthorizingUser.UserID, resp.InstallerUser.UserID)
			//backend.UpdateWorkspace(ctx, resp.TeamID, resp.TeamName)
		}
	}

	// back to the sign-up process on the main website
	if redirectURI == "" {
		c.Redirect(http.StatusTemporaryRedirect, "/start")
	} else {
		c.Redirect(http.StatusTemporaryRedirect, redirectURI)
	}
}

// UpdateAuthorization updates the authorization, or creates a new one.
func updateAuthorization(ctx context.Context, id, name, token, scope, authorizingUser, installerUser string) error {
	now := util.Timestamp()
	var auth = types.AuthorizationDS{}
	key := backend.AuthorizationKey(id)
	err := platform.DataStore().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = now
	} else {
		auth = types.AuthorizationDS{
			ID:              id,
			Name:            name,
			AccessToken:     token,
			Scope:           scope,
			AuthorizingUser: authorizingUser,
			InstallerUser:   installerUser,
			Created:         now,
			Updated:         now,
		}
	}

	_, err = platform.DataStore().Put(ctx, key, &auth)
	return err
}
