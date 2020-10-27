package api

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"

	"github.com/lnkk-app/lnkk-api/internal/stats"
)

const (
	hourlyStats = "HOURLY_STATS"
	dailyStats  = "DAILY_STATS"
)

// ScheduleHourlyTasks receives hourly cron task requests
func ScheduleHourlyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	last := platform.GetJobTimestamp(ctx, hourlyStats)

	go stats.AssetMetrics(ctx, stats.HourlyAssetMetric, "", last)
	go stats.RedirectMetrics(ctx, stats.HourlyRedirectMetric, "", last)

	platform.UpdateJob(ctx, hourlyStats, now)
	svc.StandardAPIResponse(c, nil)
}

// ScheduleDailyTasks receives daily cron task requests
func ScheduleDailyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	last := platform.GetJobTimestamp(ctx, dailyStats)

	go stats.AssetMetrics(ctx, stats.DailyAssetMetric, "", last)
	go stats.RedirectMetrics(ctx, stats.DailyRedirectMetric, "", last)

	platform.UpdateJob(ctx, dailyStats, now)
	svc.StandardAPIResponse(c, nil)
}
