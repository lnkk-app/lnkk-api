package stats

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
	"github.com/txsvc/platform/pkg/platform"
)

const (
	HourlyAssetMetric    = "HOURLY_ASSETS"
	HourlyRedirectMetric = "HOURLY_REDIRECTS"
	DailyAssetMetric     = "DAILY_ASSETS"
	DailyRedirectMetric  = "DAILY_REDIRECTS"
)

func AssetMetricsWorker(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	fmt.Println(string(bodyBytes))
	c.Status(http.StatusOK)
}

func RedirectMetricsWorker(c *gin.Context) {

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
	}
	fmt.Println(string(bodyBytes))
	c.Status(http.StatusOK)

}

// AssetMetrics runs all sorts of hourly stats in a background function
func AssetMetrics(ctx context.Context, metric, owner string, last int64) {

	n, err := urlshortener.NewAssetsSince(ctx, "", last)
	if err != nil {
		platform.ReportError(err)
		return
	}
	platform.Count(ctx, metric, owner, n)
}

// RedirectMetrics runs all sorts of hourly stats in a background function
func RedirectMetrics(ctx context.Context, metric, owner string, last int64) {

	n, err := urlshortener.RedirectsSince(ctx, "", last)
	if err != nil {
		platform.ReportError(err)
		return
	}
	platform.Count(ctx, metric, owner, n)
}
