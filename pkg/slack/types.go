package slack

type (
	// TeamInfo see https://api.slack.com/methods/team.info
	TeamInfo struct {
		OK    bool        `json:"ok"`
		Error string      `json:"error,omitempty"`
		Team  TeamElement `json:"team"`
	}

	// TeamElement see https://api.slack.com/methods/team.info
	TeamElement struct {
		ID             string       `json:"id,omitempty"`
		Name           string       `json:"name,omitempty"`
		Domain         string       `json:"domain,omitempty"`
		EmailDomain    string       `json:"email_domain,omitempty"`
		Icon           *IconElement `json:"icon,omitempty"`
		EnterpriseID   string       `json:"enterprise_id,omitempty"`
		EnterpriseName string       `json:"enterprise_name,omitempty"`
	}

	// IconElement see https://api.slack.com/methods/team.info
	IconElement struct {
		Image44      string `json:"image_34,omitempty"`
		Image68      string `json:"image_68,omitempty"`
		Image88      string `json:"image_88,omitempty"`
		Image102     string `json:"image_102,omitempty"`
		Image132     string `json:"image_132,omitempty"`
		ImageDefault bool   `json:"image_default,omitempty"`
	}
)
