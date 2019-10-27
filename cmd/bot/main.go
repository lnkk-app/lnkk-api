package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	a "github.com/lnkk-ai/lnkk/internal/api"
	"github.com/lnkk-ai/lnkk/internal/jobs"

	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
	"github.com/majordomusio/platform/pkg/tasks"
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
	router.GET(api.SchedulerBaseURL+"/workspace", a.UpdateWorkspaces)
	router.GET(api.SchedulerBaseURL+"/messages", a.CollectMessages)

	// jobs endpoints, used by the taskqueue
	router.POST(api.JobsBaseURL+"/channels", jobs.UpdateChannelsJob)
	router.POST(api.JobsBaseURL+"/users", jobs.UpdateUsersJob)

	//jobsNS.POST("/msgs", jobs.CollectMessagesJob)

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
	tasks.Close()

	log.Printf("... done")
}
