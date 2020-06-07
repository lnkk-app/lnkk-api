package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/lnkk-ai/lnkk/pkg/api"
	"github.com/lnkk-ai/lnkk/pkg/platform"
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

	// load static assets and templates
	router.Use(static.Serve("/assets/css", static.LocalFile("./public/assets/css", true)))
	router.Use(static.Serve("/assets/javascript", static.LocalFile("./public/assets/javascript", true)))
	router.LoadHTMLGlob("public/templates/*")
	// other static endpoints
	router.GET("/robots.txt", robotsEndpoint)

	// static routes
	router.GET("/", staticIndexEndpoint)
	router.GET("/error", staticErrorEndpoint)
	router.GET("/addtoslack", staticAddAppEndpoint)

	// authenticate the app
	router.GET("/a/auth", api.OAuthEndpoint)

	// API endpoints and callbacks
	router.POST("/a/actions", api.ActionRequestEndpoint)

	// start the router on port 8080, unless $PORT is set to something else
	router.Run()
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}
