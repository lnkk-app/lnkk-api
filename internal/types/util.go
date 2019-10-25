package types

import "github.com/lnkk-ai/lnkk/pkg/api"

// AsExternal create an external representation of the asset
func (t *AssetDS) AsExternal() *api.Asset {
	asset := api.Asset{
		URI:       t.URI,
		URL:       t.URL,
		SecretID:  t.SecretID,
		Cohort:    t.Cohort,
		Affiliate: t.Affiliate,
		Tags:      t.Tags,
	}
	return &asset
}
