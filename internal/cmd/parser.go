package cmd

import (
	"context"
	e "errors"
	"fmt"
	"strings"

	"github.com/txsvc/slack/pkg/slack"
)

type (
	// CommandFunc handles a command
	CommandFunc func(ctx context.Context, cmd *slack.SlashCommand, cmdLn []string) (interface{}, error)
)

var commandLookup map[string]CommandFunc = make(map[string]CommandFunc)

// ParseCommandLine pareses the command line and executes the appropriate handlers
func ParseCommandLine(ctx context.Context, group string, cmd *slack.SlashCommand) (interface{}, error) {
	parts := strings.Split(cmd.Txt, " ")
	if len(parts) == 0 {
		return nil, e.New(fmt.Sprintf("Invalid command line '%s'", cmd.Txt))
	}
	handler := commandLookup[strings.ToLower(parts[0])]
	if handler == nil {
		return nil, e.New(fmt.Sprintf("No handler for sub-command '%s'", parts[0]))
	}
	return handler(ctx, cmd, parts[1:])
}

// RegisterCmdHandler adds a slash-cmd handler
func RegisterCmdHandler(cmd, subcmd string, h CommandFunc) {
	commandLookup[strings.ToLower(subcmd)] = h
}
