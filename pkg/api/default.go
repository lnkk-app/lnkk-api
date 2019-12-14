package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// DefaultEndpoint maps to GET /
func DefaultEndpoint(c *gin.Context) {
	// TODO: real implementation, logging & auditing
	c.JSON(http.StatusOK, gin.H{"app": FullName, "vesion": Version, "status": "ok"})
}
