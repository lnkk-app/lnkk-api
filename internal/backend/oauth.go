package backend

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/appengine/memcache"

	"github.com/majordomusio/commons/pkg/errors"

	"github.com/lnkk-ai/lnkk/internal/types"
	"github.com/lnkk-ai/lnkk/pkg/platform"
)

// GetAuthorization returns the authorization granted to an app
func GetAuthorization(ctx context.Context, id string) (*types.AuthorizationDS, error) {
	var auth = types.AuthorizationDS{}
	key := "workspace.auth" + id
	_, err := memcache.Gob.Get(ctx, key, &auth)

	if err != nil {
		err = platform.DataStore().Get(ctx, AuthorizationKey(id), &auth)
		if err == nil {
			cache := memcache.Item{}
			cache.Key = key
			cache.Object = auth
			cache.Expiration, _ = time.ParseDuration(DefaultCacheDuration)
			memcache.Gob.Set(ctx, &cache)
		} else {
			return nil, err
		}
	}

	return &auth, nil
}

// GetAuthToken returns the oauth token of the workspace integration
func GetAuthToken(ctx context.Context, id string) (string, error) {
	auth, err := GetAuthorization(ctx, id)
	if err != nil {
		return "", err
	}
	if auth == nil {
		return "", errors.New(fmt.Sprintf("No authorization token for workspace '%s'", id))
	}
	return auth.AccessToken, nil
}
