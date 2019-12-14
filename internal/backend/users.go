package backend

import (
	"context"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/majordomusio/platform/pkg/store"

	"github.com/lnkk-ai/lnkk/internal/types"
)

// UpdateUser updates the user metadata
func UpdateUser(ctx context.Context, id, team, name, realName, firstName, lastName, email string, deleted, bot bool) error {
	now := util.Timestamp()
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
		user.Updated = now
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
			Created:   now,
			Updated:   now,
		}
	}

	_, err = store.Client().Put(ctx, key, &user)
	return err
}
