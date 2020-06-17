package types

type (
	// Asset is the basic entity used in the shortener
	Asset struct {
		// URI is the unique identifier for the asset
		URI string `json:"uri,omitempty"`
		// URL is the assets real url
		URL string `json:"url" binding:"required"`
		// Owner identifies who the asset belongs to
		Owner string `json:"owner,omitempty"`
		// SecretID allows admin access to this asset
		SecretID string `json:"secret_id,omitempty"`
		// Source indicates where the link originated from
		Source string `json:"source,omitempty"`
		// Client app the asset belongs to
		Client string `json:"client,omitempty"`
		// Affiliate references any affiliation
		Affiliate string `json:"affiliate,omitempty"`
		// Tags holds a comma separated list of tags
		Tags string `json:"tags,omitempty"`
	}
)
