package slack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Timestamp returns the seconds part of a Slack timestamp
// Example: "1533028651.000211" -> 1533028651
func Timestamp(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	return i
}

// TimestampNano returns a Slack timestamp as nanoseconds
func TimestampNano(ts string) int64 {
	s := strings.Split(ts, ".")
	i, _ := strconv.ParseInt(s[0], 10, 64)
	j, _ := strconv.ParseInt(s[1], 10, 64)

	return (i * 1000000) + j
}

// TimestampNanoString converts a Slack TS in nanoseconds into Slack's string representation
func TimestampNanoString(ts int64) string {
	_p1 := ts / 1000000
	_p2 := ts % 1000000
	return fmt.Sprintf("%d.%06d", _p1, _p2)
}

// SlashCommandPayload extracts the payload from a POST received by a slash command
func SlashCommandPayload(c *gin.Context) map[string]string {
	p := make(map[string]string)

	p["team_id"] = c.PostForm("team_id")
	p["team_domain"] = c.PostForm("team_domain")
	p["enterprise_id"] = c.PostForm("enterprise_id")
	p["enterprise_name"] = c.PostForm("enterprise_name")
	p["channel_id"] = c.PostForm("channel_id")
	p["channel_name"] = c.PostForm("channel_name")
	p["user_id"] = c.PostForm("user_id")
	p["user_name"] = c.PostForm("user_name")
	p["command"] = c.PostForm("command")
	p["text"] = c.PostForm("text")
	p["response_url"] = c.PostForm("response_url")
	p["trigger_id"] = c.PostForm("trigger_id")

	return p
}
