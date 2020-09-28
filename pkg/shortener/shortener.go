package shortener

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/commons/pkg/util"

	"github.com/lnkk-app/lnkk-api/pkg/types"
	"github.com/txsvc/platform/pkg/platform"
)

// CreateAsset stores an asset in the Datastore
func CreateAsset(ctx context.Context, as *types.Asset) (*types.Asset, error) {

	uri, _ := util.ShortUUID()
	secret, _ := util.ShortUUID()

	ownerType := strings.ToLower(as.Owner.Type)
	if ownerType == "" {
		ownerType = "session"
	} else if ownerType != "user" && ownerType != "session" {
		err := errors.New(fmt.Sprintf("Unknown owner type '%s'", ownerType))
		platform.Report(err)
		return nil, err
	}

	asset := AssetDS{
		LongLink:    as.LongLink,
		ShortLink:   uri,
		PreviewLink: uri,
		Owner:       as.Owner.ID,
		OwnerType:   ownerType,
		SecretID:    as.Owner.SecretID,
		Secret:      secret,
		Parent:      as.Parent,
		Created:     util.Timestamp(),
	}

	k := assetKey(uri)
	if _, err := platform.DataStore().Put(ctx, k, &asset); err != nil {
		platform.Report(err)
		return nil, err
	}

	// update the request body
	as.ShortLink = uri
	as.PreviewLink = uri
	as.Secret = secret

	// return the enriched asset instance
	return as, nil
}

// GetAsset retrieves the asset
func GetAsset(ctx context.Context, uri string) (*types.Asset, error) {
	var as AssetDS
	k := assetKey(uri)

	if err := platform.DataStore().Get(ctx, k, &as); err != nil {
		return nil, err
	}

	return as.AsExternal(), nil
}

// CreateMeasurement records a link activation
func CreateMeasurement(ctx context.Context, m *MeasurementDS) error {

	// anonimize the IP to be GDPR compliant
	m.IP = anonimizeIP(m.IP)

	// TODO: use a queue here, go routine will not work !
	CreateGeoLocation(ctx, m.IP)

	k := datastore.IncompleteKey(DatastoreMeasurement, nil)
	if _, err := platform.DataStore().Put(ctx, k, m); err != nil {
		platform.Report(err)
		return err
	}

	return nil
}

// CreateGeoLocation looks up the IP's geolocation if it is unknown
func CreateGeoLocation(ctx context.Context, ip string) error {

	var loc GeoLocationDS
	k := geoLocationKey(ip)

	if err := platform.DataStore().Get(ctx, k, &loc); err != nil {
		// assuming the location is unknown
		l, err := lookupGeoLocation(ip)
		if err != nil {
			platform.Report(err)
			return err
		}

		if _, err := platform.DataStore().Put(ctx, k, l.AsInternal()); err != nil {
			platform.Report(err)
			return err
		}
	}

	return nil
}

// Anonimize the IP to be GDPR compliant
func anonimizeIP(ip string) string {
	if strings.ContainsRune(ip, 58) {
		// IPv6
		parts := strings.Split(ip, ":")
		return fmt.Sprintf("%s:%s:%s:0000:0000:0000:0000:0000", parts[0], parts[1], parts[2])
	}
	// IPv4
	parts := strings.Split(ip, ".")
	return fmt.Sprintf("%s.%s.%s.0", parts[0], parts[1], parts[2])
}

func assetKey(uri string) *datastore.Key {
	return datastore.NameKey(DatastoreAssets, uri, nil)
}

func geoLocationKey(ip string) *datastore.Key {
	return datastore.NameKey(DatastoreGeoLocation, ip, nil)
}
