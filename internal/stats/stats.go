package stats

import (
	"context"

	"github.com/lnkk-app/lnkk-api/internal/urlshortener"
	"github.com/txsvc/platform/pkg/platform"
)

const (
	HourlyAssetMetric    = "HOURLY_ASSETS"
	HourlyRedirectMetric = "HOURLY_REDIRECTS"
	DailyAssetMetric     = "DAILY_ASSETS"
	DailyRedirectMetric  = "DAILY_REDIRECTS"
)

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
