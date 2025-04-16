package storage

import (
	"encoding/json"
	"os"

	nats "github.com/nats-io/nats.go"
)

type NatsStorage struct {
	db      *nats.Conn
	lastErr error
}

// NewNatsStorage Creates a new NATS storage
func NewNatsStorage() *NatsStorage {

	ds := &NatsStorage{}
	opts := []nats.Option{nats.Name("wasmvision NATS storage")}

	natsURL := os.Getenv("WASMVISION_STORAGE_NATS_URL")
	if natsURL == "" {
		natsURL = nats.DefaultURL
	}

	natsUserCredentialsFile := os.Getenv("WASMVISION_STORAGE_NATS_USER_CREDENTIALS")
	if natsUserCredentialsFile != "" {
		opts = append(opts, nats.UserCredentials(natsUserCredentialsFile))
	}

	natsNKeyFile := os.Getenv("WASMVISION_STORAGE_NATS_NKEY_FILE")
	if natsNKeyFile != "" {
		opt, err := nats.NkeyOptionFromSeed(natsNKeyFile)
		if err != nil {
			ds.lastErr = err
			return ds
		}
		opts = append(opts, opt)
	}

	natsTLSCert := os.Getenv("WASMVISION_STORAGE_NATS_CERT_FILE")
	natsTLSKey := os.Getenv("WASMVISION_STORAGE_NATS_KEY_FILE")

	if natsTLSCert != "" && natsTLSKey != "" {
		opts = append(opts, nats.ClientCert(natsTLSCert, natsTLSKey))
	}

	natsTLSCACert := os.Getenv("WASMVISION_STORAGE_NATS_CA_CERT")
	if natsTLSCACert != "" {
		opts = append(opts, nats.RootCAs(natsTLSCACert))
	}

	db, err := nats.Connect(natsURL, opts...)
	if err != nil {
		ds.lastErr = err
		return ds
	}

	ds.db = db

	return ds
}

// Err returns last operational error
func (ds *NatsStorage) Err() error {
	return ds.lastErr
}

// Get returns a data value from the store.
func (ds *NatsStorage) Get(root string, key string) (string, bool) {
	return "", false
}

// GetKeys returns all the keys for a specific id from the store.
func (ds *NatsStorage) GetKeys(root string) ([]string, bool) {

	return nil, false
}

// Set sets a config value in the store.
func (ds *NatsStorage) Set(root string, key string, val string) error {

	m := map[string]string{
		"proc":  root,
		"key":   key,
		"value": val,
	}

	b, err := json.Marshal(&m)
	if err != nil {
		ds.lastErr = err
		return err
	}

	err = ds.db.Publish(root, b)
	if err != nil {
		ds.lastErr = err
		return err
	}
	ds.lastErr = ds.db.Flush()

	return ds.lastErr
}

// Delete deletes data from the store.
func (ds *NatsStorage) Delete(root string, key string) {}

// DeleteAll deletes all data for a specific id from the store.
func (ds *NatsStorage) DeleteAll(root string) {
}

// Exists returns true if there is any data for a specific id in the store.
func (ds *NatsStorage) Exists(root string) bool {
	return true
}

// Close closes Redis Storage
func (ds *NatsStorage) Close() error {
	ds.db.Close()
	return nil
}
