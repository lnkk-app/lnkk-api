package backend

import (
	"fmt"

	"cloud.google.com/go/datastore"
)

// AuthorizationKey creates a datastore key for a workspace authorization based on the team_id.
func AuthorizationKey(id string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, id, nil)
}

// WorkspaceKey creates a datastore key for a workspace based on the team_id.
func WorkspaceKey(id string) *datastore.Key {
	return datastore.NameKey(DatastoreWorkspaces, id, nil)
}

// UserKey creates a datastore key for users based on the id and team_id.
func UserKey(id, team string) *datastore.Key {
	return datastore.NameKey(DatastoreUsers, team+"."+id, nil)
}

// ChannelKey creates a datastore key for channels based on the id and team_id.
func ChannelKey(id, team string) *datastore.Key {
	return datastore.NameKey(DatastoreChannels, team+"."+id, nil)
}

// MessageKey creates a datastore key for messages based on the id, timestamp and user.
func MessageKey(id, ts, user string) *datastore.Key {
	return datastore.NameKey(DatastoreMessages, MessageKeyString(id, ts, user), nil)
}

// MessageKeyString returns the string used to generate the message key
func MessageKeyString(id, ts, user string) string {
	return ts + "." + id + "." + user
}

// AttachmentKey creates a datastore key for attachements based on the message id and attachement id
func AttachmentKey(msg string, id int) *datastore.Key {
	return datastore.NameKey(DatastoreAttachments, AttachmentKeyString(msg, id), nil)
}

// AttachmentKeyString returns the string used to generate the attachment key
func AttachmentKeyString(msg string, id int) string {
	return fmt.Sprintf("%s.%d", msg, id)
}

// ReactionKey creates a datastore key for attachements based on the message id and attachement id
func ReactionKey(msg string, id int) *datastore.Key {
	return datastore.NameKey(DatastoreAttachments, ReactionKeyString(msg, id), nil)
}

// ReactionKeyString returns the string used to generate the attachment key
func ReactionKeyString(msg string, id int) string {
	return fmt.Sprintf("%s.%d", msg, id)
}
