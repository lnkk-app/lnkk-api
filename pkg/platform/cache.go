package platform

import (
	"context"
	"time"

	"google.golang.org/appengine/memcache"
)

// GetBlob returns an object from the cache
func GetBlob(ctx context.Context, key string, blob interface{}) error {
	_, err := memcache.Gob.Get(ctx, key, &blob)
	return err
}

// SetBlob adds an object to the cache
func SetBlob(ctx context.Context, key, duration string, blob interface{}) error {
	cache := memcache.Item{}
	cache.Key = key
	cache.Object = blob
	cache.Expiration, _ = time.ParseDuration(duration)
	return memcache.Gob.Set(ctx, &cache)
}

// Get returns a string
func Get(ctx context.Context, key string) (string, error) {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return string(item.Value), nil
}

// Set stores a string
func Set(ctx context.Context, key, value, duration string) error {
	exp, _ := time.ParseDuration(duration)
	cache := &memcache.Item{
		Key:        key,
		Value:      []byte(value),
		Expiration: exp,
	}
	return memcache.Set(ctx, cache)
}
