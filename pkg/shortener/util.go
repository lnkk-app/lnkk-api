package shortener

import "github.com/lnkk-app/lnkk-api/pkg/types"

// AsExternal create an external representation of the asset
func (t *AssetDS) AsExternal() *types.Asset {
	asset := types.Asset{
		LongLink:    t.LongLink,
		ShortLink:   t.ShortLink,
		PreviewLink: t.PreviewLink,
		Owner: types.AssetOwner{
			ID:       t.Owner,
			Type:     t.OwnerType,
			SecretID: t.SecretID,
		},
		Secret: t.Secret,
		Parent: t.Parent,
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
