package datastore

import "github.com/wasmvision/wasmvision/datastore/storage"

// Frames is the store for data that is specific to a single frame.
// It is used to store data that is associated with a specific frame in a video or image sequence.
// The data is stored in a map where the key is the frame number and the value is another map
// that contains key-value pairs of data associated with that frame.
type Frames struct {
	store *storage.MemStorage[int]
}

// NewFrames creates a new Frames data store.
func NewFrames() *Frames {
	return &Frames{
		store: storage.NewMemStorage[int](),
	}
}

// Get returns a data value for a specific frame from the store .
func (s *Frames) Get(frame int, key string) (string, bool) {
	return s.store.Get(frame, key)
}

// GetKeys returns all the keys for a specific frame from the store.
func (s *Frames) GetKeys(frame int) ([]string, bool) {
	return s.store.GetKeys(frame)
}

// Set sets a key/value for a specific frame in the store.
func (s *Frames) Set(frame int, key, val string) error {
	return s.store.Set(frame, key, val)
}

// Delete deletes data for a specific frame from the store.
func (s *Frames) Delete(frame int, key string) {
	s.store.Delete(frame, key)
}

// DeleteAll deletes all data for a specific frame from the store.
func (s *Frames) DeleteAll(frame int) {
	s.store.DeleteAll(frame)
}

// Exists returns true if there is any data for a specific frame in the store.
func (s *Frames) Exists(frame int) bool {
	return s.store.Exists(frame)
}
