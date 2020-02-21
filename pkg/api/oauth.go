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

		// FIXME remove this
		// LOG log.Printf("oauth: %v\n\n", util.PrintJSON(resp))

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
			platform.Report(slack.NewSimpleError("team.info", teamInfo.Error))
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		err = updateAuthorization(ctx, teamInfo.Team.ID, teamInfo.Team.Name, resp.AccessToken, resp.TokenType, resp.Scope, resp.AppID, resp.BotUserID)
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

// UpdateAuthorization updates the authorization, or creates a new one.
func updateAuthorization(ctx context.Context, id, name, token, tokenType, scope, appID, botID string) error {
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
			ID:          id,
			Name:        name,
			AccessToken: token,
			TokenType:   tokenType,
			Scope:       scope,
			AppID:       appID,
			BotUserID:   botID,
			Created:     now,
			Updated:     now,
		}
	}

	_, err = platform.DataStore().Put(ctx, key, &auth)
	return err
}
