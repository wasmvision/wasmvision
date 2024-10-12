package config

// Store is a store for config info.
type Store struct {
	storeMap map[string]string
}

// NewStore creates a new Store.
func NewStore() *Store {
	return &Store{
		storeMap: map[string]string{},
	}
}

// Get returns a config value from the store.
func (s *Store) Get(key string) (string, bool) {
	val, ok := s.storeMap[key]
	return val, ok
}

// Set sets a config value in the store.
func (s *Store) Set(key, val string) error {
	s.storeMap[key] = val
	return nil
}

// Delete deletes a frame from the cache.
func (s *Store) Delete(key string) {
	delete(s.storeMap, key)
}
