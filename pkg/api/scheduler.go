package api

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/appengine"

	"github.com/txsvc/commons/pkg/util"
	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"
)

// ScheduleHourlyTasks receives hourly cron task requests
func ScheduleHourlyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	platform.UpdateJob(ctx, "HOURLY", now)

	svc.StandardAPIResponse(c, nil)
}

// ScheduleDailyTasks receives daily cron task requests
func ScheduleDailyTasks(c *gin.Context) {
	ctx := appengine.NewContext(c.Request)

	now := util.Timestamp()
	platform.UpdateJob(ctx, "DAILY", now)

	svc.StandardAPIResponse(c, nil)
}
