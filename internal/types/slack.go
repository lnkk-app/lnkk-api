package types

type (
	// Authorization holds basic information about a Slack team/workspace and
	// the OAuth given at installtion time.
	Authorization struct {
		ID              string
		Name            string
		AccessToken     string
		Scope           string
		AuthorizingUser string
		InstallerUser   string
		Created         int64
		Updated         int64
	}
)
