package storage

import (
	"errors"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
)

type DatastoreBoltDB struct {
	db      *bolt.DB
	lastErr error
}

// NewBoltDBStorage Creates a new BoltDB storage
func NewBoltDBStorage() *DatastoreBoltDB {

	ds := &DatastoreBoltDB{}

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

	return &DatastoreBoltDB{
		db: db,
	}
}

// Error returns last operational error
func (ds *DatastoreBoltDB) Error() error {
	return ds.lastErr
}

// Get returns a data value from the store.
func (ds *DatastoreBoltDB) Get(root string, key string) (string, bool) {

	value := ""

	err := ds.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(root))
		if b == nil {
			return errors.New("invalid bucket")
		}
		value = string(b.Get([]byte(key)))
		return nil
	})

	ds.lastErr = err

	return value, value != ""
}

// GetKeys returns all the keys for a specific id from the store.
func (ds *DatastoreBoltDB) GetKeys(root string) ([]string, bool) {

	var keys []string

	err := ds.db.View(func(tx *bolt.Tx) error {
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

	ds.lastErr = err

	return keys, len(keys) > 0
}

// Set sets a config value in the store.
func (ds *DatastoreBoltDB) Set(root string, key string, val string) error {

	err := ds.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(val))
	})

	ds.lastErr = err

	return err
}

// Delete deletes data from the store.
func (ds *DatastoreBoltDB) Delete(root string, key string) {

	err := ds.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(root))
		if err != nil {
			return err
		}
		return b.Delete([]byte(key))
	})

	ds.lastErr = err
}

// DeleteAll deletes all data for a specific id from the store.
func (ds *DatastoreBoltDB) DeleteAll(root string) {

	err := ds.db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(root))
	})

	ds.lastErr = err
}

// Exists returns true if there is any data for a specific id in the store.
func (ds *DatastoreBoltDB) Exists(root string) bool {
	exists := false

	err := ds.db.View(func(tx *bolt.Tx) error {
		exists = tx.Bucket([]byte(root)) != nil
		return nil
	})

	ds.lastErr = err

	return exists

}

// Close closes BoltDB Storage
func (ds *DatastoreBoltDB) Close() error {
	return ds.db.Close()
}
