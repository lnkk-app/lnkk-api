package platform

// Implements a simple key-value store (strings only!), backed by the Google Datastore, with a simplistic caching mechanism.

import (
	"context"

	"cloud.google.com/go/datastore"
	"github.com/golang/groupcache/lru"
	"github.com/majordomusio/commons/pkg/util"
)

const (
	// DatastoreKV collection KV
	DatastoreKV string = "KV"
	// DefaultExpiration is the minimum time to keep an entry in the store
	DefaultExpiration int64 = 86400 * 14 // 14 days
)

type (
	// KV is the datastructure to store stuff
	KV struct {
		Key     string
		Value   string
		Expires int64
	}
)

var cache = lru.New(100)

// Set adds an entry to the store. Existing values get updated
func Set(ctx context.Context, k, v string, duration int64) error {
	key := datastore.NameKey(DatastoreKV, k, nil)
	var kv = KV{
		Key:     k,
		Value:   v,
		Expires: util.Timestamp() + duration,
	}

	_, err := DataStore().Put(ctx, key, &kv)
	if err != nil {
		return err
	}

	// add also to the LRU cache
	cache.Add(&key, &kv)
	return nil
}

// Get retrieves a value from the store or raises an exception if it does not exist
func Get(ctx context.Context, k string) (string, error) {
	key := datastore.NameKey(DatastoreKV, k, nil)

	// check if the value is in the cache ...
	v, ok := cache.Get(&key)
	if ok {
		return v.(KV).Value, nil
	}

	var kv = KV{}
	err := DataStore().Get(ctx, key, &kv)
	if err == nil {
		return kv.Value, nil
	}

	return "", err
}

// Invalidate removes an entry from the cache
func Invalidate(ctx context.Context, k string) {
	key := datastore.NameKey(DatastoreKV, k, nil)
	cache.Remove(&key)
	DataStore().Delete(ctx, key)
}
