package types

type (
	// AuthorizationDS holds basic information about a Slack team/workspace and
	// the OAuth given at installtion time.
	AuthorizationDS struct {
		ID              string
		Name            string
		AccessToken     string
		Scope           string
		AuthorizingUser string
		InstallerUser   string
		// internal
		Created int64
		Updated int64
	}
)
