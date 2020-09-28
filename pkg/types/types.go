package types

type (
	// Asset is the basic entity used in the shortener
	Asset struct {
		LongLink    string     `json:"long_link" binding:"required"`
		ShortLink   string     `json:"short_link,omitempty"`
		PreviewLink string     `json:"preview_link,omitempty"`
		Suffix      string     `json:"suffix,omitempty"` // SHORT or UNGUESSABLE
		Secret      string     `json:"secret,omitempty"`
		Owner       AssetOwner `json:"owner"`
		Parent      string     `json:"parent,omitempty"`
	}

	// AssetOwner describes the asset ownership
	AssetOwner struct {
		ID       string `json:"id" binding:"required"`
		Type     string `json:"type,omitempty"` // USER or SESSION
		SecretID string `json:"secret_id,omitempty"`
	}
)
