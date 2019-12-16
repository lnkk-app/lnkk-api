package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	a "github.com/majordomusio/commons/pkg/api"

	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
	"github.com/majordomusio/platform/pkg/tasks"

	"github.com/lnkk-ai/lnkk/pkg/api"
)

func main() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()

	// basic config
	gin.DisableConsoleColor()
	// a new router
	router := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// default endpoints that are not part of the API namespace
	//router.GET("/", a.DefaultEndpoint)
	router.GET("/robots.txt", a.RobotsEndpoint)

	// authenticate the app
	router.GET("/auth", api.AuthEndpoint)

	// scheduler endpoints
	router.GET(api.SchedulerBaseURL+"/workspace", scheduleUpdateWorkspaces)
	router.GET(api.SchedulerBaseURL+"/messages", scheduleCollectMessages)
	// statistics etc
	router.GET(api.SchedulerBaseURL+"/hourly", scheduleHourlyTasks)
	router.GET(api.SchedulerBaseURL+"/daily", scheduleDailyTasks)

	// jobs endpoints, used by the taskqueue
	router.POST(api.JobsBaseURL+"/j/channels", taskUpdateChannels)
	router.POST(api.JobsBaseURL+"/j/users", taskUpdateUsers)
	router.POST(api.JobsBaseURL+"/j/messages", taskCollectMessages)
	// hourly, periodic tasks
	router.POST(api.JobsBaseURL+"/j/activations", taskLinkActivations)
	router.POST(api.JobsBaseURL+"/j/daily", taskDaily)

	// start the router on port 8080, unless $PORT is set to something else
	router.Run()
}

func shutdown() {
	store.Close()
	errorreporting.Close()
	tasks.Close()

	log.Printf("Exiting ...")
}
