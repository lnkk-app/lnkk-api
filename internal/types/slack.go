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
		Created         int64
		Updated         int64
	}

	// WorkspaceDS basic workspace metadata
	WorkspaceDS struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// Scheduler
		NextUpdate     int64 `json:"-"` // scheduled time of the next crawl
		UpdateSchedule int   `json:"-"` // interval, in seconds
		// Internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// User data
	User struct {
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
		NextCrawl       int64 `json:"-"` // scheduled time of the next crawl
		CrawlerSchedule int   `json:"-"` // interval, in seconds
		// Internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}

	// Message holds slack messages
	Message struct {
		ChannelID       string
		TeamID          string
		User            string
		TS              int64  // ts the message was created, according to the Slack API
		Text            string `datastore:",noindex"`
		HasAttachements bool
		// Internal
		Created int64 // ts this record was created
		Updated int64 // ts this record was last updated
	}

	// Attachment related to messages
	Attachment struct {
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

	// ArchivedMessage is the body of a message
	ArchivedMessage struct {
		Text           string                      `json:"text,omitempty"`
		User           string                      `json:"user,omitempty"`
		Created        string                      `json:"created,omitempty"`
		Timestamp      int64                       `json:"ts,omitempty"`
		HasAttachments bool                        `json:"has_attachments"`
		Attachments    []ArchivedMessageAttachment `json:"attachments,omitempty"`
	}

	// ArchivedMessageAttachment holds the metadata of an attachment
	ArchivedMessageAttachment struct {
		Text         string `json:"text,omitempty"`
		FallbackText string `json:"fallback_text,omitempty"`
	}
)
