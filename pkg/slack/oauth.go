package slack

import (
	"encoding/json"
	e "errors"
	"net/http"
	"os"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/lnkk-ai/lnkk/internal/backend"
	"github.com/lnkk-ai/lnkk/pkg/errors"
	"github.com/lnkk-ai/lnkk/pkg/platform"
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
		resp, err := GetOAuthToken(ctx, code)

		if err != nil {
			platform.Report(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error")
			return
		}

		// get team info
		var teamInfo TeamInfo
		err = Get(ctx, resp.AccessToken, "team.info", "", &teamInfo)
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

// GetOAuthToken exchanges a temporary OAuth verifier code for an access token
func GetOAuthToken(ctx context.Context, code string) (*OAuthResponse, error) {

	url := SlackEndpoint + "oauth.v2.access?code=" + code

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(os.Getenv(SlackClientID), os.Getenv(SlackClientSecret))

	// post the request to Slack
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// unmarshal the response
	var response OAuthResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	return &response, err
}
