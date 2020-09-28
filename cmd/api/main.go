package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"

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

	// API endpoints and callbacks
	apiNamespace := r.Group(api.APIPrefix)
	apiNamespace.POST("/short", api.ShortenEndpoint)
	// redirect endpoint
	r.GET("/r/:uri", api.RedirectEndpoint)

	return r
}

func staticIndexEndpoint(c *gin.Context) {
	svc.StandardAPIResponse(c, nil)
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}
