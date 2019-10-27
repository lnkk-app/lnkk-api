package api

const (
	// FullName is the name of the service
	FullName string = "lnkk.social"
	// Version is the human readable version string of this build
	Version string = "1.0"
)

const (
	// BotBaseURL is the prefix for all public API endpoints
	BotBaseURL string = "/a/1"
	// SchedulerBaseURL is the prefix for all scheduller/cron tasks
	SchedulerBaseURL string = "/_i/1/scheduler"
	// JobsBaseURL is the prefix for all scheduled jobs
	JobsBaseURL string = "/_i/1/jobs"
)

type (
	// Asset is the basic entity used in the shortener
	Asset struct {
		// URI is the unique identifier for the asset
		URI string `json:"uri,omitempty"`
		// URL is the assets real url
		URL string `json:"url" binding:"required"`
		// SecretID allows admin access
		SecretID string `json:"secret_id,omitempty"`
		// Source indicates where the link will be used
		Source string `json:"source,omitempty"`
		// Cohort the asset belongs to
		Cohort string `json:"cohort,omitempty"`
		// Affiliate references any affiliation
		Affiliate string `json:"affiliate,omitempty"`
		// Tags holds a comma separated list of tags
		Tags string `json:"tags,omitempty"`
	}
)
