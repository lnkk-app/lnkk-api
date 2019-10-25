package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// standardResponse is the default way to respond to API requests
func standardResponse(c *gin.Context, err error) {
	if err == nil {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	} else {
		//platform.Error(platform.NewRequestContext(c), "api.response", err.Error())
		// TODO proper error handling. For now 400 it is
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}

// standardJSONResponse is the default way to respond to API requests
func standardJSONResponse(c *gin.Context, res interface{}, err error) {
	if err == nil {
		if res == nil {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			c.JSON(http.StatusOK, res)
		}
	} else {
		//platform.Error(platform.NewRequestContext(c), "api.response", err.Error())
		// TODO proper error handling. For now 400 it is
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "msg": err.Error()})
	}
}
