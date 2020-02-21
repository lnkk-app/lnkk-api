package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/majordomusio/commons/pkg/env"
)

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

// robotsEndpoint maps to GET /robots.txt
func robotsEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /a/\n"))
}
