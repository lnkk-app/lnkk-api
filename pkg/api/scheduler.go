package api

import (
	"fmt"

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

	platform.CreateSimpleTask(ctx, WorkerBaseURL+"/metrics/assets", fmt.Sprintf("%s:%s:%d", stats.HourlyAssetMetric, "-", last))
	platform.CreateSimpleTask(ctx, WorkerBaseURL+"/metrics/redirects", fmt.Sprintf("%s:%s:%d", stats.HourlyRedirectMetric, "-", last))

	platform.UpdateJob(ctx, hourlyStats, now)
	svc.StandardAPIResponse(c, nil)
}

// ScheduleDailyTasks receives daily cron task requests
func ScheduleDailyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	last := platform.GetJobTimestamp(ctx, dailyStats)

	platform.CreateSimpleTask(ctx, WorkerBaseURL+"/metrics/assets", fmt.Sprintf("%s:%s:%d", stats.DailyAssetMetric, "-", last))
	platform.CreateSimpleTask(ctx, WorkerBaseURL+"/metrics/redirects", fmt.Sprintf("%s:%s:%d", stats.DailyRedirectMetric, "-", last))

	platform.UpdateJob(ctx, dailyStats, now)
	svc.StandardAPIResponse(c, nil)
}
