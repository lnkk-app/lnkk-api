package store

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/datastore"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/util"
)

var dsClient *datastore.Client

func init() {
	ctx := context.Background()
	c, err := datastore.NewClient(ctx, env.Getenv("PROJECT_ID", ""))
	if err != nil {
		errorreporting.Report(err)
	}
	dsClient = c
}

// Close does the clean-up
func Close() {
	dsClient.Close()
}

// AssetKey creates the datastore key for an asset
func AssetKey(uri string) *datastore.Key {
	return datastore.NameKey(types.DatastoreAssets, uri, nil)
}

// GeoLocationKey creates the datastore key for a geolocation
func GeoLocationKey(ip string) *datastore.Key {
	return datastore.NameKey(types.DatastoreGeoLocation, ip, nil)
}

// CreateAsset stores an asset in the Datastore
func CreateAsset(ctx context.Context, as *types.AssetDS) error {
	as.Created = util.Timestamp()

	k := AssetKey(as.URI)
	if _, err := dsClient.Put(ctx, k, as); err != nil {
		errorreporting.Report(err)
		return err
	}

	return nil
}

// GetAsset retrieves the asset
func GetAsset(ctx context.Context, uri string) (*api.Asset, error) {
	var as types.AssetDS
	k := AssetKey(uri)

	if err := dsClient.Get(ctx, k, &as); err != nil {
		return nil, err
	}

	return as.AsExternal(), nil
}

// CreateMeasurement records a link activation
func CreateMeasurement(ctx context.Context, m *types.MeasurementDS) error {

	// anonimize the IP to be GDPR compliant
	m.IP = anonimizeIP(m.IP)

	// TODO: use a queue here, go routine will not work !
	CreateGeoLocation(ctx, m.IP)

	k := datastore.IncompleteKey(types.DatastoreMeasurement, nil)
	if _, err := dsClient.Put(ctx, k, m); err != nil {
		errorreporting.Report(err)
		return err
	}

	return nil
}

// CreateGeoLocation looks up the IP's geolocation if it is unknown
func CreateGeoLocation(ctx context.Context, ip string) error {

	var loc types.GeoLocationDS
	k := GeoLocationKey(ip)

	if err := dsClient.Get(ctx, k, &loc); err != nil {
		// assuming the location is unknown
		l, err := lookupGeoLocation(ip)
		if err != nil {
			errorreporting.Report(err)
			return err
		}

		if _, err := dsClient.Put(ctx, k, l.AsInternal()); err != nil {
			errorreporting.Report(err)
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
