package types

type (
	// AuthorizationDS holds basic information about a Slack team/workspace and
	// the OAuth given at installtion time.
	AuthorizationDS struct {
		ID          string
		Name        string
		AccessToken string
		TokenType   string
		AppID       string
		BotUserID   string
		Scope       string
		// internal
		Created int64
		Updated int64
	}
)
