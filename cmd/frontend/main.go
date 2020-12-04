package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/txsvc/platform/pkg/platform"
	"github.com/txsvc/service/pkg/svc"
)

func init() {
	// setup shutdown handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		shutdown()
		os.Exit(1)
	}()
}

func shutdown() {
	platform.Close()
	log.Printf("Exiting ...")
}

func main() {

	// add basic routes
	svc.AddDefaultEndpoints()
	svc.ServeStaticAssets("/", "./cmd/frontend/public")

	// add custom endpoints with authentication
	//api := svc.SecureGroup("/api", a.MiddlewareFunc())
	//api.GET("/public", "chat.read", testAPIResponse)
	//api.POST("/private", "chat.write", testAPIResponse)

	// add CORS handler, allowing all. See https://github.com/gin-contrib/cors
	//svc.Use(cors.Default())

	// add session handler
	//store := cookie.NewStore([]byte(secret))
	//svc.Use(sessions.Sessions("svcexample", store))

	// add the service/router to a server on $PORT and launch it. This call BLOCKS !
	svc.Start()
}
