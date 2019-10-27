package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	a "github.com/lnkk-ai/lnkk/internal/api"

	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/errorreporting"
	"github.com/lnkk-ai/lnkk/pkg/store"
)

const (
	// BotBaseURL is the prefix for all public API endpoints
	BotBaseURL string = "/a/1"
	// SchedulerBaseURL is the prefix for all scheduller/cron tasks
	SchedulerBaseURL string = "/_i/1/scheduler"
	// JobsBaseURL is the prefix for all scheduled jobs
	JobsBaseURL string = "/_i/1/jobs"
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
	router.GET("/", defaultEndpoint)
	router.GET("/robots.txt", a.RobotsEndpoint)

	// authenticate the app
	router.GET("/auth", a.AuthEndpoint)

	// scheduler endpoints
	router.GET(SchedulerBaseURL+"/workspace", a.UpdateWorkspaces)
	router.GET(SchedulerBaseURL+"/messages", a.CollectMessages)

	// // group of internal endpoints for scheduling
	//scheduleNS := router.Group(types.SchedulerBaseURL)
	//scheduleNS.GET("/workspace", scheduler.UpdateWorkspaces)
	//scheduleNS.GET("/msgs", scheduler.CollectMessages)

	// start the router on port 8080, unless ENV PORT is set to something else
	router.Run()
}

// defaultEndpoint maps to GET /*
func defaultEndpoint(c *gin.Context) {
	// TODO: real implementation, logging & auditing
	c.JSON(http.StatusOK, gin.H{"app": "bot", "vesion": api.Version, "status": "ok"})
}

func shutdown() {
	log.Printf("Shutting down ...")

	store.Close()
	errorreporting.Close()

	log.Printf("... done")
}
