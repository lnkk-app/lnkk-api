package backend

import (
	"context"

	"github.com/majordomusio/commons/pkg/util"

	"github.com/majordomusio/platform/pkg/store"

	"github.com/lnkk-ai/lnkk/internal/types"
)

// UpdateAuthorization updates the authorization, or creates a new one.
func UpdateAuthorization(ctx context.Context, id, name, token, scope, authorizingUser, installerUser string) error {
	now := util.Timestamp()
	var auth = types.AuthorizationDS{}
	key := AuthorizationKey(id)
	err := store.Client().Get(ctx, key, &auth)

	if err == nil {
		auth.AccessToken = token
		auth.Scope = scope
		auth.Updated = now
	} else {
		auth = types.AuthorizationDS{
			ID:              id,
			Name:            name,
			AccessToken:     token,
			Scope:           scope,
			AuthorizingUser: authorizingUser,
			InstallerUser:   installerUser,
			Created:         now,
			Updated:         now,
		}
	}

	_, err = store.Client().Put(ctx, key, &auth)
	return err
}
