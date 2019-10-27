package slack

const (
	// SlackEndpoint is the URI of the Slack API
	SlackEndpoint string = "https://slack.com/api/"
	// SlackClientID is the Client ID of App
	SlackClientID string = "SLACK_CLIENT_ID"
	// SlackClientSecret is the Client Secret of the App
	SlackClientSecret string = "SLACK_CLIENT_SECRET"
	// SlackOAuthToken is the default OAuth token
	SlackOAuthToken string = "SLACK_OAUTH_TOKEN"
	// SlackVerificationToken is a secret token used to verify requests from Slack
	SlackVerificationToken string = "SLACK_VERIFICATION_TOKEN"
	// SlackResponseTypeChannel is used to send messages to channels that are visible to evryone
	SlackResponseTypeChannel string = "in_channel"
	// SlackResponseTypeEphemeral is used to send a message to a channel that is only visible to the current user
	SlackResponseTypeEphemeral string = "ephemeral"
)

type (

	// Notification is a short message sent to a channel via a webhook
	Notification struct {
		Username string `json:"username,omitempty"`
		Channel  string `json:"channel,omitempty"`
		Text     string `json:"text"`
		Pretext  string `json:"pretext,omitempty"`
	}

	//
	// Standard Messages
	//

	// StandardMessage is used to give a simple reponse to the user as feedback to a custom action reuqest
	StandardMessage struct {
		Text         string               `json:"text,omitempty"`
		ResponseType string               `json:"response_type,omitempty"` // 'in_channel' or 'ephemeral'
		Attachements []MessageAttachement `json:"attachments,omitempty"`
	}

	// MessageAttachement lets you add more context to a message, making them more useful and effective.
	// See https://api.slack.com/docs/message-attachments
	MessageAttachement struct {
		Title     string `json:"title,omitempty"`
		TitleLink string `json:"title_link,omitempty"`
		Pretext   string `json:"pretext,omitempty"`
		Text      string `json:"text,omitempty"`
	}

	// StandardResponse is the generic response received after a Web API request.
	// See https://api.slack.com/web#responses
	StandardResponse struct {
		OK      bool   `json:"ok"`
		Stuff   string `json:"stuff,omitempty"`
		Warning string `json:"warning,omitempty"`
		Error   string `json:"error,omitempty"`
	}

	//
	// Custom Message Action
	//

	// ActionRequest is the payload received from Slack when the user triggers a custom message action
	ActionRequest struct {
		Type             string               `json:"type,omitempty"`
		Token            string               `json:"token,omitempty"`
		ActionTimestamp  string               `json:"action_ts,omitempty"`
		Team             MessageActionTeam    `json:"team,omitempty"`
		User             MessageActionUser    `json:"user,omitempty"`
		Channel          MessageActionChannel `json:"channel,omitempty"`
		CallbackID       string               `json:"callback_id,omitempty"`
		TriggerID        string               `json:"trigger_id,omitempty"`
		MessageTimestamp string               `json:"message_ts,omitempty"`
		Message          ActionRequestMessage `json:"message,omitempty"`
		ResponseURL      string               `json:"response_url,omitempty"`
		Submission       map[string]string    `json:"submission,omitempty"`
	}

	// ActionRequestMessage is the message's main content
	ActionRequestMessage struct {
		Type         string                     `json:"type,omitempty"`
		User         string                     `json:"user,omitempty"`
		Text         string                     `json:"text,omitempty"`
		Attachements []ActionRequestAttachement `json:"attachments,omitempty"`
		Timestamp    string                     `json:"ts,omitempty"`
	}

	// ActionRequestAttachement describes message attachements such as links or files
	ActionRequestAttachement struct {
		ServiceName string `json:"service_name,omitempty"`
		Title       string `json:"title,omitempty"`
		TitleLink   string `json:"title_link,omitempty"`
		Text        string `json:"text,omitempty"`
		Fallback    string `json:"fallback,omitempty"`
		ImageURL    string `json:"image_url,omitempty"`
		FromURL     string `json:"from_url,omitempty"`
		ImageWidth  int    `json:"image_width,omitempty"`
		ImageHeight int    `json:"image_height,omitempty"`
		ImageBytes  int    `json:"image_bytes,omitempty"`
		ServiceIcon string `json:"service_icon,omitempty"`
		ID          int    `json:"id,omitempty"`
		OriginalURL string `json:"original_url,omitempty"`
	}

	// MessageActionTeam identifies the Slack workspace the message originates from
	MessageActionTeam struct {
		ID     string `json:"id,omitempty"`
		Domain string `json:"domain,omitempty"`
	}

	// MessageActionUser identifies the user who triggered the custom action
	MessageActionUser struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	// MessageActionChannel identifies the channel the custom action was triggered from
	MessageActionChannel struct {
		ID   string `json:"id,omitempty"`
		Name string `json:"name,omitempty"`
	}

	//
	// Dialog
	//

	// OpenDialogRequest is used to request a dialog being openended in the Slack client
	OpenDialogRequest struct {
		TriggerID string `json:"trigger_id"`
		Dialog    Dialog `json:"dialog"`
	}

	// Dialog defines the dialog and its fields
	Dialog struct {
		Title          string          `json:"title,omitempty"`
		CallbackID     string          `json:"callback_id,omitempty"`
		NotifyOnCancel bool            `json:"notify_on_cancel,omitempty"`
		Elements       []DialogElement `json:"elements,omitempty"`
	}

	// DialogElement describes a field in the dialog
	DialogElement struct {
		Type        string `json:"type,omitempty"`
		Subtype     string `json:"subtype,omitempty"`
		Label       string `json:"label,omitempty"`
		Name        string `json:"name,omitempty"`
		Placeholder string `json:"placeholder,omitempty"`
	}

	//
	// Channels
	//

	// ChannelsInfoResponse returns a channel object.
	ChannelsInfoResponse struct {
		StandardResponse
		Channel Channel `json:"channel,omitempty"`
	}

	// ChannelsListResponse returns a list of paginated channel objects.
	ChannelsListResponse struct {
		StandardResponse
		Channels         []Channel         `json:"channels,omitempty"`
		ResponseMetadata map[string]string `json:"response_metadata,omitempty"`
	}

	// ChannelsHistoryResponse is a list of message events from a specified public channel.
	ChannelsHistoryResponse struct {
		StandardResponse
		Latest   string           `json:"latest,omitempty"`
		Messages []ChannelMessage `json:"messages,omitempty"`
		HasMore  bool             `json:"has_more,omitempty"`
	}

	// Channel  contains information about a workspace channel
	// See https://api.slack.com/types/channel
	Channel struct {
		ID                 string             `json:"id,omitempty"`
		Name               string             `json:"name,omitempty"`
		IsChannel          bool               `json:"is_channel,omitempty"`
		Created            int64              `json:"created,omitempty"`
		Creator            string             `json:"creator,omitempty"`
		IsArchived         bool               `json:"is_archived,omitempty"`
		IsGeneral          bool               `json:"is_general,omitempty"`
		NameNormalized     string             `json:"name_normalized,omitempty"`
		IsShared           bool               `json:"is_shared,omitempty"`
		IsOrgShared        bool               `json:"is_org_shared,omitempty"`
		IsMember           bool               `json:"is_member,omitempty"`
		IsPrivate          bool               `json:"is_private,omitempty"`
		IsMPIM             bool               `json:"is_mpim,omitempty"`
		LastRead           string             `json:"last_read,omitempty"`
		Latest             ChannelMessage     `json:"latest,omitempty"`
		UnreadCount        int                `json:"unread_count,omitempty"`
		UnreadCountDisplay int                `json:"unread_count_display,omitempty"`
		Members            []string           `json:"members,omitempty"`
		Topic              ChannelDescription `json:"topic,omitempty"`
		Purpose            ChannelDescription `json:"purpose,omitempty"`
		PreviousNames      []string           `json:"previous_names,omitempty"`
	}

	// ChannelMessage is a message posted on a channel.
	ChannelMessage struct {
		Text         string               `json:"text,omitempty"`
		Username     string               `json:"username,omitempty"`
		User         string               `json:"user,omitempty"`
		BotID        string               `json:"bot_id,omitempty"`
		Attachements []ChannelAttachement `json:"attachments,omitempty"`
		Reactions    []ChannelReaction    `json:"reactions,omitempty"`
		Type         string               `json:"type,omitempty"`
		Subtype      string               `json:"subtype,omitempty"`
		IsStarred    bool                 `json:"is_starred,omitempty"`
		TS           string               `json:"ts,omitempty"`
	}

	// ChannelDescription is some metadata about the channels purpose etc.
	ChannelDescription struct {
		Value   string `json:"value,omitempty"`
		Creator string `json:"creator,omitempty"`
		LastSet int64  `json:"last_set,omitempty"`
	}

	// ChannelAttachement are attachments on a message in the channel
	ChannelAttachement struct {
		Text     string `json:"text,omitempty"`
		ID       int    `json:"id,omitempty"`
		Fallback string `json:"fallback,omitempty"`
	}

	// ChannelReaction list a reaction and all users that used it.
	ChannelReaction struct {
		Name  string   `json:"name,omitempty"`
		Count int      `json:"count,omitempty"`
		Users []string `json:"users,omitempty"`
	}

	//
	// Groups (private channels)
	//

	// GroupsListResponse is a list of group objects (also known as "private channel objects")
	GroupsListResponse struct {
		StandardResponse
		Groups []Group `json:"groups,omitempty"`
	}

	// Group holds information about private channels (groups)
	Group struct {
		ID         string             `json:"id,omitempty"`
		Name       string             `json:"name,omitempty"`
		Created    int64              `json:"created,omitempty"`
		Creator    string             `json:"creator,omitempty"`
		IsArchived bool               `json:"is_archived,omitempty"`
		Members    []string           `json:"members,omitempty"`
		Topic      ChannelDescription `json:"topic,omitempty"`
		Purpose    ChannelDescription `json:"purpose,omitempty"`
	}

	//
	// Users
	//

	// UsersInfoResponse returns a user object
	UsersInfoResponse struct {
		StandardResponse
		User User `json:"user,omitempty"`
	}

	// UsersListResponse returns a list of paginated user objects, in no particular order
	UsersListResponse struct {
		StandardResponse
		Members          []User            `json:"members,omitempty"`
		CacheTS          int               `json:"cache_ts,omitempty"`
		ResponseMetadata map[string]string `json:"response_metadata,omitempty"`
	}

	// User contains information about a member
	// See https://api.slack.com/types/user
	User struct {
		ID                string  `json:"id,omitempty"`
		TeamID            string  `json:"team_id,omitempty"`
		Name              string  `json:"name,omitempty"`
		Deleted           bool    `json:"deleted,omitempty"`
		Color             string  `json:"color,omitempty"`
		RealName          string  `json:"real_name,omitempty"`
		TZ                string  `json:"tz,omitempty"`
		TimezoneLabel     string  `json:"tz_label,omitempty"`
		TimezoneOffset    int     `json:"tz_offset,omitempty"`
		Profile           Profile `json:"profile,omitempty"`
		IsAdmin           bool    `json:"is_admin,omitempty"`
		IsOwner           bool    `json:"is_owner,omitempty"`
		IsPrimaryOwner    bool    `json:"is_primary_owner,omitempty"`
		IsRestricted      bool    `json:"is_restricted,omitempty"`
		IsUltraRestricted bool    `json:"is_ultra_restricted,omitempty"`
		IsBot             bool    `json:"is_bot,omitempty"`
		Updated           int     `json:"updated,omitempty"`
		IsAppUser         bool    `json:"is_app_user,omitempty"`
		Has2FA            bool    `json:"has_2fa,omitempty"`
	}

	// Profile object contains as much information as the user has supplied in the default profile
	Profile struct {
		AvatarHash            string `json:"avatar_hash,omitempty"`
		StatusText            string `json:"status_text,omitempty"`
		StatusEmoji           string `json:"status_emoji,omitempty"`
		FirstName             string `json:"first_name,omitempty"`
		LastName              string `json:"last_name,omitempty"`
		RealName              string `json:"real_name,omitempty"`
		DisplayName           string `json:"display_name,omitempty"`
		RealNameNormalized    string `json:"real_name_normalized,omitempty"`
		DisplayNameNormalized string `json:"display_name_normalized,omitempty"`
		Email                 string `json:"email,omitempty"`
		Image24               string `json:"image_24,omitempty"`
		Image32               string `json:"image_32,omitempty"`
		Image72               string `json:"image_72,omitempty"`
		Image192              string `json:"image_192,omitempty"`
		Image512              string `json:"image_512,omitempty"`
		Team                  string `json:"team,omitempty"`
	}
)
