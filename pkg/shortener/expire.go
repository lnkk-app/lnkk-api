package shortener

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/env"
	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/lnkk-app/lnkk-api/internal/misc"
)

// AssetExpirationWorker receives worker tasks to expire asset
func AssetExpirationWorker(c *gin.Context) {

	payload, err := misc.ExtractBodyAsString(c)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		c.Status(http.StatusOK)
		return
	}

	ts, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		platform.ReportError(errors.NewOperationError(fmt.Sprintf("Invalid parameter. Expected string(int64), got: %d of type %T", ts, ts), nil))
		c.Status(http.StatusOK)
		return
	}

	ctx := appengine.NewContext(c.Request)
	cutOffTimestamp := ts - (env.GetInt("MAX_ASSET_AGE", ExpireAfter) * 86400)

	n, err := ExpireAssets(ctx, cutOffTimestamp)
	if err != nil {
		platform.ReportError(err)
	} else {
		platform.Count(ctx, DailyExpiration, "", n)
	}

	// all OK
	c.Status(http.StatusOK)
}

// ExpireAssets looks for expired assets and changes their state
func ExpireAssets(ctx context.Context, ts int64) (int, error) {
	var q *datastore.Query
	count := 0
	now := util.Timestamp()

	q = datastore.NewQuery(DatastoreAssets).Filter("LastAccess <", ts).Filter("State =", StateActive)
	it := platform.DataStore().Run(ctx, q)
	for {
		var asset Asset

		_, err := it.Next(&asset)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return -1, err
		}

		// expire the asset
		asset.State = StateExpired
		asset.Modified = now
		k := assetKey(asset.ShortLink)
		if _, err := platform.DataStore().Put(ctx, k, &asset); err != nil {
			return -1, err
		}

		count++
	}
	return count, nil
}
