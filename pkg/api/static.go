package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RobotsEndpoint maps to GET /robots.txt
func RobotsEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /a/\n"))
}

// AdsEndpoint maps to GET /ads.txt
func AdsEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /a/\n"))
}

// HumansEndpoint maps to GET /humans.txt
func HumansEndpoint(c *gin.Context) {
	// simply write text back ...
	c.Header("Content-Type", "text/plain")

	// a simple robots.txt file, disallow the API
	c.Writer.Write([]byte("User-agent: *\n\n"))
	c.Writer.Write([]byte("Disallow: /a/\n"))
}

// NullEndpoint just respondes with an empty 200
func NullEndpoint(c *gin.Context) {
	c.Status(http.StatusAccepted)
}
