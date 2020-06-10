package slack

import (
	"github.com/gin-gonic/gin"
)

// see https://api.slack.com/interactivity/slash-commands

type (
	// SlashCommand encapsulates the payload sent from Slack as a result of invoking a /slash command
	SlashCommand struct {
		TeamID         string
		TeamDomain     string
		EnterpriseID   string
		EnterpriseName string
		ChannelID      string
		ChannelName    string
		UserID         string
		UserName       string
		Command        string
		Txt            string
		ResponseURL    string
		TriggerID      string
		Token          string // DEPRECATED
	}
)

// GetSlashCommand extracts the payload from a POST received by a slash command
func GetSlashCommand(c *gin.Context) *SlashCommand {
	return &SlashCommand{
		TeamID:         c.PostForm("team_id"),
		TeamDomain:     c.PostForm("team_domain"),
		EnterpriseID:   c.PostForm("enterprise_id"),
		EnterpriseName: c.PostForm("enterprise_name"),
		ChannelID:      c.PostForm("channel_id"),
		ChannelName:    c.PostForm("channel_name"),
		UserID:         c.PostForm("user_id"),
		UserName:       c.PostForm("user_name"),
		Command:        c.PostForm("command"),
		Txt:            c.PostForm("text"),
		ResponseURL:    c.PostForm("response_url"),
		TriggerID:      c.PostForm("trigger_id"),
		Token:          c.PostForm("token"), // DEPRECATED
	}
}
