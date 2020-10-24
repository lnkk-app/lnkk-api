package api

import (
	"github.com/gin-gonic/gin"

	"github.com/txsvc/service/pkg/svc"
)

// ScheduleHourlyTasks receives hourly cron task requests
func ScheduleHourlyTasks(c *gin.Context) {
	//topic := "api.shorten.post"
	//ctx := appengine.NewContext(c.Request)

	svc.StandardAPIResponse(c, nil)
}

// ScheduleDailyTasks receives daily cron task requests
func ScheduleDailyTasks(c *gin.Context) {
	//topic := "api.shorten.post"
	//ctx := appengine.NewContext(c.Request)

	svc.StandardAPIResponse(c, nil)
}
