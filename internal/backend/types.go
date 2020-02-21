package backend

const (
	// DatastoreAuthorizations collection AUTHORIZATION
	DatastoreAuthorizations string = "AUTHORIZATIONS"

	// BackgroundWorkQueue is the default background job queue
	BackgroundWorkQueue string = "background-work"
	// DefaultCacheDuration default time to keep stuff in memory
	DefaultCacheDuration string = "15m"
	// SchedulerBaseURL is the prefix for all scheduller/cron tasks
	SchedulerBaseURL string = "/_i/1/scheduler"
	// JobsBaseURL is the prefix for all scheduled jobs
	JobsBaseURL string = "/_i/1/jobs"
)

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
