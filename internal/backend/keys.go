package backend

import (
	"cloud.google.com/go/datastore"
)

// AuthorizationKey creates a datastore key for a workspace authorization based on the team_id.
func AuthorizationKey(id string) *datastore.Key {
	return datastore.NameKey(DatastoreAuthorizations, id, nil)
}
