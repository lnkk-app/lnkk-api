package slack

import (
	"golang.org/x/net/context"
)

// UsersList returns a list of all users in the workspace. This includes deleted/deactivated users
// See https://api.slack.com/methods/users.list
func UsersList(ctx context.Context, token, cursor string) (*UsersListResponse, error) {
	cmd := "users.list"
	var q string

	if cursor != "" {
		q = "cursor=" + cursor
	}

	var resp UsersListResponse
	err := GetRequest(ctx, token, cmd, q, &resp)

	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, newSlackError(cmd, resp.Error)
	}

	return &resp, nil
}

// UsersInfo gets information about a user, identified by its user id
// See https://api.slack.com/methods/users.info
func UsersInfo(ctx context.Context, token, userID string) (*UsersInfoResponse, error) {
	cmd := "users.info"
	q := "user=" + userID

	var resp UsersInfoResponse
	err := GetRequest(ctx, token, cmd, q, &resp)

	if err != nil {
		return nil, err
	}

	if !resp.OK {
		return nil, newSlackError(cmd, resp.Error)
	}

	return &resp, nil
}
