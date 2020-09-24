package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"

	"github.com/txsvc/platform/pkg/platform"
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

	// add endpoints here ...

	// start the router on port 8080, unless $PORT is set to something else
	router.Run()
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}
