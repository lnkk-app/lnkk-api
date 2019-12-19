package slack

import (
	"encoding/json"
	"net/http"
	"os"

	"golang.org/x/net/context"
)

type (

	// OAuthResponse is used to give a simple reponse to the user as feedback to a custom action reuqest
	OAuthResponse struct {
		OK              bool         `json:"ok,omitempty"`
		AccessToken     string       `json:"access_token,omitempty"`
		TokenType       string       `json:"token_type,omitempty"`
		Scope           string       `json:"scope,omitempty"`
		AppID           string       `json:"app_id,omitempty"`
		AppUserID       string       `json:"app_user_id,omitempty"`
		TeamName        string       `json:"team_name,omitempty"`
		TeamID          string       `json:"team_id,omitempty"`
		IncomingWebhook ScopeWebhook `json:"incoming_webhook,omitempty"`
		Bot             ScopeBot     `json:"bot,omitempty"`
		AuthorizingUser OAuthUser    `json:"authorizing_user,omitempty"`
		InstallerUser   OAuthUser    `json:"app_home,omitempty"`
		Scopes          Scopes       `json:"scopes,omitempty"`
	}

	// OAuthUser is a user authorizing an app or installing an app
	OAuthUser struct {
		UserID  string `json:"user_id,omitempty"`
		AppHome string `json:"app_home,omitempty"`
	}

	// Scopes list different permissions
	Scopes struct {
		AppHome []string `json:"app_home,omitempty"`
		Team    []string `json:"team,omitempty"`
		Channel []string `json:"channel,omitempty"`
		Group   []string `json:"group,omitempty"`
		MPIM    []string `json:"mpim,omitempty"`
		IM      []string `json:"im,omitempty"`
		User    []string `json:"user,omitempty"`
		// TODO: is this all? Idon't trust the API docs...
		// See https://api.slack.com/methods/oauth.access
	}

	// ScopeWebhook not sure?
	ScopeWebhook struct {
		URL              string `json:"url,omitempty"`
		Channel          string `json:"channel,omitempty"`
		ConfigurationURL string `json:"configuration_url,omitempty"`
	}

	// ScopeBot not sure?
	ScopeBot struct {
		BotUserID      string `json:"bot_user_id,omitempty"`
		BotAccessToken string `json:"bot_access_token,omitempty"`
	}
)

// OAuthAccess exchanges a temporary OAuth verifier code for an access token
func OAuthAccess(ctx context.Context, code string) (*OAuthResponse, error) {

	url := SlackEndpoint + "oauth.access?code=" + code

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
