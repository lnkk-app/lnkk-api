package statistics

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
)

// AssetMetricsWorker receives worker tasks to create asset metrics
func AssetMetricsWorker(c *gin.Context) {

	payload := ""
	if c.Request.Body != nil {
		body, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// just report and return, resending will not change anything
			platform.ReportError(err)
			c.Status(http.StatusOK)
			return
		}
		payload = string(body)
	}

	parts := strings.Split(payload, ":")
	if len(parts) != 3 {
		// just report and return, resending will not change anything
		platform.ReportError(errors.NewOperationError(fmt.Sprintf("Invalid number of parameters. Expected 3, got %d", len(parts)), nil))
		c.Status(http.StatusOK)
		return
	}
	metric := parts[0]
	owner := parts[1]
	if owner == "-" {
		owner = ""
	}
	last, err := strconv.ParseInt(parts[2], 10, 64)
	if err == nil {
		platform.ReportError(errors.NewOperationError(fmt.Sprintf("Invalid parameter. Expected string(int64), got: %d of type %T", last, last), nil))
		c.Status(http.StatusOK)
		return
	}

	//AssetMetrics(appengine.NewContext(c.Request), parts[0], owner, last)

	ctx := appengine.NewContext(c.Request)
	n, err := AssetsSince(ctx, owner, last)
	if err != nil {
		platform.ReportError(err)
		return
	}
	platform.Count(ctx, metric, owner, n)

	// all OK
	c.Status(http.StatusOK)
}

// AssetsSince returns the number of new assets since a given point in time
func AssetsSince(ctx context.Context, owner string, ts int64) (int, error) {
	var q *datastore.Query

	if owner != "" {
		q = datastore.NewQuery(urlshortener.DatastoreAssets).Filter("Owner =", owner).Filter("Created >=", ts).KeysOnly()
	} else {
		q = datastore.NewQuery(urlshortener.DatastoreAssets).Filter("Created >=", ts).KeysOnly()
	}

	n, err := platform.DataStore().Count(ctx, q)
	if err != nil {
		return -1, err
	}
	return n, nil
}
