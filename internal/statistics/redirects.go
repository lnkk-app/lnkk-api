package statistics

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/datastore"
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/errors"
	"github.com/txsvc/platform/pkg/platform"

	"github.com/lnkk-app/lnkk-api/internal/misc"
	"github.com/lnkk-app/lnkk-api/pkg/shortener"
)

// RedirectMetricsWorker receives worker tasks to create redirect metrics
func RedirectMetricsWorker(c *gin.Context) {

	payload, err := misc.ExtractBodyAsString(c)
	if err != nil {
		// just report and return, resending will not change anything
		platform.ReportError(err)
		c.Status(http.StatusOK)
		return
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
	if err != nil {
		platform.ReportError(errors.NewOperationError(fmt.Sprintf("Invalid parameter. Expected string(int64), got: %d of type %T", last, last), nil))
		c.Status(http.StatusOK)
		return
	}

	ctx := appengine.NewContext(c.Request)
	n, err := RedirectsSince(ctx, owner, last)
	if err != nil {
		platform.ReportError(err)
	} else {
		platform.Count(ctx, metric, owner, n)
	}

	// all OK
	c.Status(http.StatusOK)

}

// RedirectsSince returns the number of asset redirects since a given point in time
func RedirectsSince(ctx context.Context, owner string, ts int64) (int, error) {
	var q *datastore.Query

	if owner != "" {
		q = datastore.NewQuery(shortener.DatastoreRedirectHistory).Filter("Owner =", owner).Filter("Created >=", ts).KeysOnly()
	} else {
		q = datastore.NewQuery(shortener.DatastoreRedirectHistory).Filter("Created >=", ts).KeysOnly()
	}

	n, err := platform.DataStore().Count(ctx, q)
	if err != nil {
		return -1, err
	}
	return n, nil
}
