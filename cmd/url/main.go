package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/platform/pkg/errorreporting"
	"github.com/majordomusio/platform/pkg/store"
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
	router := gin.New()
	router.Use(gin.Recovery())

	// endpoints
	router.POST("/cmd/s", slashCmdShortenEndpoint)
	router.GET("/r/:uri", redirectEndpoint)

	// start the router on port 8080, unless ENV[PORT] is set to something else
	router.Run()
}

func shutdown() {
	store.Close()
	errorreporting.Close()

	log.Printf("Exiting ...")
}
