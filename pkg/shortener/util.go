package shortener

import "github.com/lnkk-ai/lnkk/pkg/types"

// AsExternal create an external representation of the asset
func (t *AssetDS) AsExternal() *types.Asset {
	asset := types.Asset{
		URI:       t.URI,
		URL:       t.URL,
		Owner:     t.Owner,
		SecretID:  t.SecretID,
		Client:    t.Client,
		Affiliate: t.Affiliate,
		Tags:      t.Tags,
	}
	return &asset
}

// AsInternal converts the geolocation into the internal DS struct
func (r *LocationType) AsInternal() *GeoLocationDS {
	location := GeoLocationDS{
		IP:          r.IP,
		Host:        r.Host,
		ISP:         r.Isp,
		City:        r.City,
		CountryCode: r.Countrycode,
		CountryName: r.Countryname,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
	}
	return &location
}
