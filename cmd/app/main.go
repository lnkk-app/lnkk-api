package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
	"github.com/majordomusio/commons/pkg/services"

	s "github.com/lnkk-ai/lnkk/internal/slack"
	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"
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

	// setup slack handlers & callbacks
	slack.RegisterSlashCmdHandler("/lnkk", s.CmdLnkkHandler)
	slack.RegisterStartAction("add_newsletter", s.StartAddToNewsletter)
	slack.RegisterCompleteAction("add_newsletter", s.CompleteAddToNewsletter)

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

	// routes to load static assets and templates
	r.Use(static.Serve("/assets/css", static.LocalFile("./public/assets/css", true)))
	r.Use(static.Serve("/assets/javascript", static.LocalFile("./public/assets/javascript", true)))
	r.LoadHTMLGlob("public/templates/*")

	// default static endpoints
	r.GET("/robots.txt", services.RobotsEndpoint)
	r.GET("/ads.txt", services.NullEndpoint)    // FIXME change to the real handler
	r.GET("/humans.txt", services.NullEndpoint) // FIXME change to the real handler

	// static routes
	r.GET("/", staticIndexEndpoint)
	r.GET("/error", staticErrorEndpoint)
	r.GET("/addtoslack", staticAddAppEndpoint)

	// API endpoints and callbacks

	// Slack integration
	r.POST("/s/1/actions", slack.ActionRequestEndpoint)
	r.POST("/s/1/cmd/lnk", slack.SlashCmdEndpoint)
	r.GET("/s/1/auth", slack.OAuthEndpoint)

	// generic API
	r.POST("/a/1/short", api.ShortenEndpoint)
	r.GET("/r/:uri", api.RedirectEndpoint)

	return r
}

func staticIndexEndpoint(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

func staticErrorEndpoint(c *gin.Context) {
	c.HTML(http.StatusOK, "error.tmpl", gin.H{})
}

func staticAddAppEndpoint(c *gin.Context) {
	c.HTML(http.StatusOK, "add.tmpl", gin.H{
		"scope":     env.Getenv("SLACK_OAUTH_SCOPE", "commands,incoming-webhook,team:read"),
		"client_id": env.Getenv("SLACK_CLIENT_ID", ""),
	})
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}
