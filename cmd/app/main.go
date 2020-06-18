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

	// basic config
	gin.DisableConsoleColor()
	// a new router
	router := gin.New()
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	// routes to load static assets and templates
	router.Use(static.Serve("/assets/css", static.LocalFile("./public/assets/css", true)))
	router.Use(static.Serve("/assets/javascript", static.LocalFile("./public/assets/javascript", true)))
	router.LoadHTMLGlob("public/templates/*")

	// default static endpoints
	router.GET("/robots.txt", api.RobotsEndpoint)
	router.GET("/ads.txt", api.NullEndpoint)    // FIXME change to the real handler
	router.GET("/humans.txt", api.NullEndpoint) // FIXME change to the real handler

	// static routes
	router.GET("/", staticIndexEndpoint)
	router.GET("/error", staticErrorEndpoint)
	router.GET("/addtoslack", staticAddAppEndpoint)

	// API endpoints and callbacks

	// Slack integration
	router.POST("/a/actions", api.ActionRequestEndpoint)
	router.POST("/a/cmd/lnk", api.SlashCmdEndpoint)
	router.GET("/a/auth", slack.OAuthEndpoint)

	// generic API
	router.POST("/a/short", api.ShortenEndpoint)
	router.GET("/r/:uri", api.RedirectEndpoint)

	// start the router on port 8080, unless $PORT is set to something else
	router.Run()
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
