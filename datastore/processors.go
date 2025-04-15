package datastore

// Processors is the store for data that is not specific to a single frame but applies to a single processor.
type Processors struct {
	store DataStorage
}

// NewProcessors creates a new Processors data store.
func NewProcessors(s DataStorage) *Processors {
	return &Processors{
		store: s,
	}
}

// Get returns a data value from the store.
func (s *Processors) Get(processor string, key string) (string, bool) {
	return s.store.Get(processor, key)
}

// GetKeys returns all the keys for a specific processor from the store.
func (s *Processors) GetKeys(processor string) ([]string, bool) {
	return s.store.GetKeys(processor)
}

// Set sets a config value in the store.
func (s *Processors) Set(processor string, key string, val string) error {
	return s.store.Set(processor, key, val)
}

// Delete deletes data from the store.
func (s *Processors) Delete(processor string, key string) {
	s.store.Delete(processor, key)
}

// DeleteAll deletes all data for a specific processor from the store.
func (s *Processors) DeleteAll(processor string) {
	s.store.DeleteAll(processor)
}

// Exists returns true if there is any data for a specific processor in the store.
func (s *Processors) Exists(processor string) bool {
	return s.store.Exists(processor)
}

// Error returns last operational error if any. nil otherwise.
func (s *Processors) Error() error {
	return nil
}

// Close closes the underlying storage.
func (s *Processors) Close() error {
	return nil
}
