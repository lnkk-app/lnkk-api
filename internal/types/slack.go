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
		// Internal
		Created int64
		Updated int64
	}

	// WorkspaceDS basic workspace metadata
	WorkspaceDS struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// Scheduler
		Next           int64 `json:"-"` // scheduled time of the next crawl
		UpdateSchedule int   `json:"-"` // interval, in seconds
		// Internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// UserDS data
	UserDS struct {
		ID        string
		TeamID    string
		Name      string
		RealName  string
		FirstName string
		LastName  string
		EMail     string
		IsDeleted bool
		IsBot     bool
		// Internal
		Created int64
		Updated int64
	}

	// `datastore:",noindex" json:"feature"`

	// ChannelDS metadata
	ChannelDS struct {
		ID         string `json:"id"`
		TeamID     string `json:"team_id"`
		Name       string `json:"name"`
		Topic      string `datastore:",noindex" json:"topic"`
		Purpose    string `datastore:",noindex" json:"purpose"`
		IsArchived bool   `json:"archived"`
		IsPrivate  bool   `json:"private"`
		IsDeleted  bool   `json:"deleted"`
		// Scheduler
		Latest          int64 `json:"-"` // ts of the last message crawl
		Next            int64 `json:"-"` // scheduled time of the next crawl
		CrawlerSchedule int   `json:"-"` // interval, in seconds
		// Internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}

	// MessageDS holds slack messages
	MessageDS struct {
		ChannelID    string
		TeamID       string
		User         string
		TS           int64  // ts the message was created, according to the Slack API
		Text         string `datastore:",noindex"`
		Attachements bool
		Reactions    bool
		// Internal
		Created int64 // ts this record was created
		Updated int64 // ts this record was last updated
	}

	// AttachmentDS related to messages
	AttachmentDS struct {
		MessageID string
		ChannelID string
		TeamID    string
		ID        int
		Text      string `datastore:",noindex"`
		Fallback  string `datastore:",noindex"`
		// Internal
		Created int64 // ts this record was created
		Updated int64 // ts this record was last updated
	}

	// ReactionDS related to messages
	ReactionDS struct {
		MessageID string
		ChannelID string
		TeamID    string
		Reaction  string
		Count     int
		Users     []string
		// Internal
		Created int64 // ts this record was created
		Updated int64 // ts this record was last updated
	}
)
