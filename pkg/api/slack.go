package api

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lnkk-ai/lnkk/pkg/platform"
	"github.com/lnkk-ai/lnkk/pkg/slack"
)

// ActionRequestEndpoint receives callbacks from Slack
func ActionRequestEndpoint(c *gin.Context) {
	var action slack.ActionRequest

	//payload := c.Request.FormValue("payload")
	err := json.Unmarshal([]byte(c.Request.FormValue("payload")), &action)
	if err != nil {
		platform.Report(err)
	}

	log.Printf("%v", action)
}
