package storage

import (
	"errors"
	"os"
	"strconv"

	bolt "go.etcd.io/bbolt"
)

type BoltDBStorage struct {
	db      *bolt.DB
	lastErr error
}

// NewBoltDBStorage Creates a new BoltDB storage
func NewBoltDBStorage() *BoltDBStorage {

	ds := &BoltDBStorage{}

	fname := os.Getenv("WASMVISION_STORAGE_BOLTDB_FILENAME")
	if fname == "" {
		ds.lastErr = errors.New("WASMVISION_STORAGE_BOLTDB_FILENAME env var not defined")
		return ds
	}

	fmodeInt := 0644
	fmode := os.Getenv("WASMVISION_STORAGE_BOLTDB_FILEMODE")
	if fmode != "" {
		i, err := strconv.ParseInt(fmode, 8, 32)
		if err == nil {
			fmodeInt = int(i)
		}
	}

	db, err := bolt.Open(fname, os.FileMode(fmodeInt), nil)
	if err != nil {
		ds.lastErr = err
		return ds
	}

	return &BoltDBStorage{
		db: db,
	}
}

// Error returns last operational error
func (ds *BoltDBStorage) Err() error {
	return ds.lastErr
}

// Get returns a data value from the store.
func (ds *BoltDBStorage) Get(root string, key string) (string, bool) {

	value := ""

	ds.lastErr = ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(root))
		if b == nil {
			return errors.New("invalid bucket")
		}
		value = string(b.Get([]byte(key)))
		return nil
	})

	return value, value != ""
}

// GetKeys returns all the keys for a specific id from the store.
func (ds *BoltDBStorage) GetKeys(root string) ([]string, bool) {

	var keys []string

	ds.lastErr = ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(root))
		if b == nil {
			return errors.New("invalid bucket")
		}
		err := b.ForEach(func(k, v []byte) error {
			keys = append(keys, string(k))
			return nil
		})
		return err
	})

	return keys, len(keys) > 0
}

// Set sets a config value in the store.
func (ds *BoltDBStorage) Set(root string, key string, val string) error {

	ds.lastErr = ds.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(val))
	})

	return ds.lastErr
}

// Delete deletes data from the store.
func (ds *BoltDBStorage) Delete(root string, key string) {

	ds.lastErr = ds.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})
}

// DeleteAll deletes all data for a specific id from the store.
func (ds *BoltDBStorage) DeleteAll(root string) {

	ds.lastErr = ds.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(root))
	})
}

// Exists returns true if there is any data for a specific id in the store.
func (ds *BoltDBStorage) Exists(root string) bool {
	exists := false

	ds.lastErr = ds.db.View(func(tx *bolt.Tx) error {
		exists = tx.Bucket([]byte(root)) != nil
		return nil
	})

	return exists
}

// Close closes BoltDB Storage
func (ds *BoltDBStorage) Close() error {
	return ds.db.Close()
}
