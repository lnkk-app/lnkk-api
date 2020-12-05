package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/slack/pkg/slack"

	"github.com/lnkk-app/lnkk-api/internal/actions"
	"github.com/lnkk-app/lnkk-api/internal/cmd"
	"github.com/lnkk-app/lnkk-api/pkg/api"
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

	// setup slack commands and actions
	setupSlackCommands()
	setupSlackActions()

	// basic http stack config
	gin.DisableConsoleColor()
	// create the routes
	router := setupRoutes()
	// start the router on port 8080, unless $PORT is set to something else
	router.Run()
}

func setupRoutes() *gin.Engine {
	// a new router
	r := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// static routes

	// api endpoints and callbacks
	apiNamespace := r.Group(api.APIPrefix)
	apiNamespace.POST("/short", api.ShortenEndpoint)
	apiNamespace.GET("/auth", slack.OAuthEndpoint)
	apiNamespace.POST("/slack/cmd", slack.SlashCmdEndpoint)
	apiNamespace.POST("/slack/action", slack.ActionRequestEndpoint)

	// hack
	// apiNamespace.GET("/migrate", api.MigrationEndpoint001)
	// end hack

	// scheduler
	schedulerNamespace := r.Group(api.SchedulerBaseURL)
	schedulerNamespace.GET("/hourly", api.ScheduleHourlyTasks)
	schedulerNamespace.GET("/daily", api.ScheduleDailyTasks)

	// redirect endpoint
	r.GET("/r/:short", api.RedirectEndpoint)

	return r
}

func setupSlackCommands() {
	slack.RegisterSlashCmdHandler("/lnkk", cmd.SlackCmdLnkk)
}

func setupSlackActions() {
	slack.RegisterStartAction("add_newsletter", actions.StartAddToNewsletter)
	slack.RegisterCompleteAction("add_newsletter", actions.CompleteAddToNewsletter)
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}
