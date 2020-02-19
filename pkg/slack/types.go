package slack

type (

	// MessageActionPayload receives
	MessageActionPayload struct {
		Type        string `json:"type"`
		CallbackID  string `json:"callback_id"`
		TriggerID   string `json:"trigger_id"`
		ResponseURL string `json:"response_url"`
		Token       string `json:"token"`
	}
)
