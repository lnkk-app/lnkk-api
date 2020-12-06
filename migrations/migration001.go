package migrations

import (
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
)

// MigrationEndpoint001 added the LastAccess attribute and set the asset state correctly
func MigrationEndpoint001(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	count := 0
	now := util.Timestamp()
	it := platform.DataStore().Run(ctx, datastore.NewQuery(urlshortener.DatastoreAssets))
	for {
		var asset urlshortener.Asset

		_, err := it.Next(&asset)
		if err == iterator.Done {
			break
		}
		if err != nil {
			platform.ReportError(err)
		}

		// update the asset
		asset.LastAccess = now
		asset.State = urlshortener.StateActive
		k := assetKey(asset.ShortLink)
		if _, err := platform.DataStore().Put(ctx, k, &asset); err != nil {
			platform.ReportError(err)
		}

		count++
	}
	fmt.Printf("Migrated %d assets.", count)
	svc.StandardAPIResponse(c, nil)
}

func assetKey(uri string) *datastore.Key {
	return datastore.NameKey(urlshortener.DatastoreAssets, uri, nil)
}
