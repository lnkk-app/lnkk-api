package backend

import (
	"context"
	"time"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/majordomusio/commons/pkg/util"
	"github.com/majordomusio/platform/pkg/store"
	"google.golang.org/appengine/memcache"
)

// UpdateUser updates the user metadata
func UpdateUser(ctx context.Context, id, team, name, realName, firstName, lastName, email string, deleted, bot bool) error {

	var user = types.UserDS{}
	key := UserKey(id, team)
	err := store.Client().Get(ctx, key, &user)

	if err == nil {
		user.Name = name
		user.RealName = realName
		user.FirstName = firstName
		user.LastName = lastName
		user.EMail = email
		user.IsDeleted = deleted
		user.IsBot = bot
		user.Updated = util.Timestamp()
	} else {
		user = types.UserDS{
			ID:        id,
			TeamID:    team,
			Name:      name,
			RealName:  realName,
			FirstName: firstName,
			LastName:  lastName,
			EMail:     email,
			IsDeleted: deleted,
			IsBot:     bot,
			Created:   util.Timestamp(),
			Updated:   util.Timestamp(),
		}
	}

	_, err = store.Client().Put(ctx, key, &user)
	return err
}

// GetUserName returns user's full name
func GetUserName(ctx context.Context, userID, teamID string) string {

	k := userID + "." + teamID
	n, err := memcache.Get(ctx, k)

	if err != nil {
		var user = types.UserDS{}
		key := UserKey(userID, teamID)
		err := store.Client().Get(ctx, key, &user)

		if err == nil {
			// add the user to the cache
			var n memcache.Item
			n.Key = k
			n.Value = ([]byte)(user.RealName)
			n.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Set(ctx, &n)

			return user.RealName
		}
		// user was not found in the datastore either
		return ""
	}
	return (string)(n.Value)
}
