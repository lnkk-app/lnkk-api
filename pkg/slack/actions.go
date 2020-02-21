package slack

type (

	// message_action -> ActionRequest
	// view_submission -> ViewSubmission

	// ActionRequestPeek is used to determin the type of request
	ActionRequestPeek struct {
		Type string `json:"type,omitempty"`
	}

	// See https://api.slack.com/reference/interaction-payloads/actions

	// ActionRequest is the payload received from Slack when the user triggers a custom message action
	// type == message_action
	ActionRequest struct {
		Type             string                `json:"type,omitempty"`
		Token            string                `json:"token,omitempty"`
		ActionTimestamp  string                `json:"action_ts,omitempty"`
		Team             *MessageActionTeam    `json:"team,omitempty"`
		User             *MessageActionUser    `json:"user,omitempty"`
		Channel          *MessageActionChannel `json:"channel,omitempty"`
		CallbackID       string                `json:"callback_id,omitempty"`
		TriggerID        string                `json:"trigger_id,omitempty"`
		MessageTimestamp string                `json:"message_ts,omitempty"`
		Message          *ActionRequestMessage `json:"message,omitempty"`
		ResponseURL      string                `json:"response_url,omitempty"`
		Submission       map[string]string     `json:"submission,omitempty"`
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

	// ViewSubmission see https://api.slack.com/reference/interaction-payloads/views#view_submission
	// type == view_submission
	ViewSubmission struct {
		Type      string             `json:"type,omitempty"`
		Team      *MessageActionTeam `json:"team,omitempty"`
		User      *MessageActionUser `json:"user,omitempty"`
		Token     string             `json:"token,omitempty"`
		TriggerID string             `json:"trigger_id,omitempty"`
		View      *ViewElement       `json:"view,omitempty"`
	}
)
