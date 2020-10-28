package urlshortener

import (
	"context"

	"cloud.google.com/go/datastore"

	"github.com/txsvc/platform/pkg/platform"
)

// NewAssetsSince returns the number of new assets since a given point in time
func NewAssetsSince(ctx context.Context, owner string, ts int64) (int, error) {
	var q *datastore.Query

	if owner != "" {
		q = datastore.NewQuery(DatastoreAssets).Filter("Owner =", owner).Filter("Created >=", ts).KeysOnly()
	} else {
		q = datastore.NewQuery(DatastoreAssets).Filter("Created >=", ts).KeysOnly()
	}

	n, err := platform.DataStore().Count(ctx, q)
	if err != nil {
		return -1, err
	}
	return n, nil
}

// RedirectsSince returns the number of new assets since a given point in time
func RedirectsSince(ctx context.Context, owner string, ts int64) (int, error) {
	var q *datastore.Query

	if owner != "" {
		q = datastore.NewQuery(DatastoreRedirectHistory).Filter("Owner =", owner).Filter("Created >=", ts).KeysOnly()
	} else {
		q = datastore.NewQuery(DatastoreRedirectHistory).Filter("Created >=", ts).KeysOnly()
	}

	n, err := platform.DataStore().Count(ctx, q)
	if err != nil {
		return -1, err
	}
	return n, nil
}
