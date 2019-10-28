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

	// WorkspaceDS basic workspace metadata
	WorkspaceDS struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		// Scheduler
		Next     int64 `json:"-"` // scheduled time of the next crawl
		Schedule int   `json:"-"` // interval, in seconds
		// internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
	}

	// UserDS data
	UserDS struct {
		ID        string `json:"id"`
		TeamID    string `json:"team_id"`
		Name      string `json:"name"`
		RealName  string `json:"real_name"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		EMail     string `json:"email"`
		IsDeleted bool   `json:"deleted"`
		IsBot     bool   `json:"bot"`
		// internal
		Created int64 `json:"-"`
		Updated int64 `json:"-"`
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
		// scheduler
		Latest   int64 `json:"-"` // ts of the last message crawl
		Next     int64 `json:"-"` // scheduled time of the next crawl
		Schedule int   `json:"-"` // interval, in seconds
		// internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}

	// MessageDS holds slack messages
	MessageDS struct {
		ChannelID string `json:"id"`
		TeamID    string `json:"team_id"`
		User      string `json:"user"`
		Text      string `datastore:",noindex" json:"text"`
		Timestamp int64  `json:"ts"` // ts the message was created, according to the Slack API
		// supporting analytics
		Type         string `json:"type"`
		Subtype      string `json:"subtype"`
		Attachements bool   `json:"attachements"`
		Reactions    bool   `json:"reactions"`
		Day          int    `json:"day_of_week"` // day of the week (Sun = 0)
		Hour         int    `json:"hour_of_day"` // hour of the day (0..23)
		// internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}

	// AttachmentDS related to messages
	AttachmentDS struct {
		MessageID   string `json:"message_id"` // it's parent record
		ChannelID   string `json:"channel_id"`
		TeamID      string `json:"team_id"`
		Index       int    `json:"index"`
		Text        string `datastore:",noindex" json:"text"`
		Alternative string `datastore:",noindex" json:"alt_text"`
		// internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}

	// ReactionDS related to messages
	ReactionDS struct {
		MessageID string   `json:"message_id"` // it's parent record
		ChannelID string   `json:"channel_id"`
		TeamID    string   `json:"team_id"`
		Reaction  string   `json:"reaction"`
		Count     int      `json:"count"`
		Users     []string `datastore:",noindex" json:"users"`
		// internal
		Created int64 `json:"-"` // ts this record was created
		Updated int64 `json:"-"` // ts this record was last updated
	}
)
