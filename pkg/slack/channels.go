package slack

import (
	"fmt"

	"golang.org/x/net/context"
)

// ChannelsList lists all channels in a Slack team
// See https://api.slack.com/methods/channels.list
func ChannelsList(ctx context.Context, token, cursor string) (*ChannelsListResponse, error) {
	cmd := "channels.list"
	var q string

	if cursor != "" {
		q = "cursor=" + cursor
	}

	var resp ChannelsListResponse
	err := GetRequest(ctx, token, cmd, q, &resp)

	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, newError(cmd, resp.Error)
	}

	return &resp, nil
}

// ChannelsInfo gets information about a channel
// See https://api.slack.com/methods/channels.info
func ChannelsInfo(ctx context.Context, token, channelID string) (*ChannelsInfoResponse, error) {
	cmd := "channels.info"
	q := "channel=" + channelID

	var resp ChannelsInfoResponse
	err := GetRequest(ctx, token, cmd, q, &resp)

	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, newError(cmd, resp.Error)
	}

	return &resp, nil
}

// ChannelsHistory fetches history of messages and events from a channel.
// See https://api.slack.com/methods/channels.history
//
// In case a timeframe is specified, the first parameter represents the
// newest message, the second the oldest message to include.
func ChannelsHistory(ctx context.Context, token, channelID string, count int, t ...string) (*ChannelsHistoryResponse, error) {
	cmd := "channels.history"
	var q string

	if len(t) == 0 {
		q = fmt.Sprintf("channel=%s&count=%d", channelID, count)
	} else if len(t) == 1 {
		q = fmt.Sprintf("channel=%s&count=%d&latest=%s", channelID, count, t[0])
	} else {
		q = fmt.Sprintf("channel=%s&count=%d&latest=%s&oldest=%s", channelID, count, t[0], t[1])
	}

	var resp ChannelsHistoryResponse
	err := GetRequest(ctx, token, cmd, q, &resp)

	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, newError(cmd, resp.Error)
	}

	return &resp, nil
}
