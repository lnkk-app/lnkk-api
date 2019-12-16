package slack

import (
	"log"

	"github.com/majordomusio/platform/pkg/errorreporting"

	"golang.org/x/net/context"
)

// see https://api.slack.com/methods/chat.postMessage

// PostSimpleMessage sends a string to a channel
func PostSimpleMessage(ctx context.Context, token, channel, msg string) {
	m := SimpleMessage{
		Channel: channel,
		Text:    msg,
	}

	r, err := PostRequest(ctx, token, "chat.postMessage", &m)
	if err != nil {
		errorreporting.Report(err)
	}
	log.Println(r)

}
