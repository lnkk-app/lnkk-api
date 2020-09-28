package api

// asExternal create an external representation of the asset
func (t *AssetDS) asExternal() *Asset {
	asset := Asset{
		LongLink:    t.LongLink,
		ShortLink:   t.ShortLink,
		PreviewLink: t.PreviewLink,
		Owner: AssetOwner{
			ID:       t.Owner,
			Type:     t.OwnerType,
			SecretID: t.SecretID,
		},
		Secret: t.Secret,
		Parent: t.Parent,
	}
	return &asset
}

// asInternal converts the geolocation into the internal DS struct
func (r *LocationType) asInternal() *GeoLocationDS {
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
