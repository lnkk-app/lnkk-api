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

/*

{
    "ok": true,
    "access_token": "xoxb-941724765634-944826281235-n75RuOTOTYh3bXPl2Cbj3a7s",
    "token_type": "bot",
    "scope": "commands,incoming-webhook,team:read",
    "app_id": "AU5V984BZ",
    "bot_user_id": "UTSQA896X",
    "incoming_webhook": {
        "url": "https://hooks.slack.com/services/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",
        "channel": "#general",
        "configuration_url": "https://lnkk-ai-dev.slack.com/services/BUALCC6DN"
    }
}

*/
