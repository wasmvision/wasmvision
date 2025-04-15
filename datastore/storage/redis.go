package storage

import (
	"context"
	"errors"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	db      *redis.Client
	lastErr error
}

// NewRedisStorage Creates a new Redis storage
func NewRedisStorage() *RedisStorage {

	ds := &RedisStorage{}

	redisURL := os.Getenv("WASMVISION_STORAGE_REDIS_URL")
	if redisURL == "" {
		ds.lastErr = errors.New("WASMVISION_STORAGE_REDIS_URL env var not defined")
		return ds
	}

	options, err := redis.ParseURL(redisURL)
	if err != nil {
		ds.lastErr = err
		return ds
	}

	ds.db = redis.NewClient(options)

	return ds
}

// Error returns last operational error
func (ds *RedisStorage) Error() error {
	return ds.lastErr
}

// Get returns a data value from the store.
func (ds *RedisStorage) Get(root string, key string) (string, bool) {

	r := ds.db.HGet(context.Background(), root, key)

	ds.lastErr = r.Err()
	return r.Val(), r.Val() != ""
}

// GetKeys returns all the keys for a specific id from the store.
func (ds *RedisStorage) GetKeys(root string) ([]string, bool) {

	r := ds.db.HKeys(context.Background(), root)

	ds.lastErr = r.Err()

	return r.Val(), len(r.Val()) != 0
}

// Set sets a config value in the store.
func (ds *RedisStorage) Set(root string, key string, val string) error {

	r := ds.db.HSet(context.Background(), root, key, val)

	ds.lastErr = r.Err()

	return ds.lastErr
}

// Delete deletes data from the store.
func (ds *RedisStorage) Delete(root string, key string) {

	r := ds.db.HDel(context.Background(), root, key)
	ds.lastErr = r.Err()
}

// DeleteAll deletes all data for a specific id from the store.
func (ds *RedisStorage) DeleteAll(root string) {
	r1 := ds.db.HGetAll(context.Background(), root)
	ds.lastErr = r1.Err()

	if ds.lastErr != nil {
		return
	}

	keys := make([]string, 0, len(r1.Val()))

	for k := range r1.Val() {
		keys = append(keys, k)
	}

	r2 := ds.db.HDel(context.Background(), root, keys...)
	ds.lastErr = r2.Err()
}

// Exists returns true if there is any data for a specific id in the store.
func (ds *RedisStorage) Exists(root string) bool {

	r := ds.db.HLen(context.Background(), root)

	return r.Val() > 0
}

// Close closes Redis Storage
func (ds *RedisStorage) Close() error {
	return ds.db.Close()
}
