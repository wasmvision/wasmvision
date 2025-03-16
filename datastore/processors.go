package datastore

type Processors struct {
	storeMap map[string]map[string]string
}

// NewProcessors creates a new Processors data store.
func NewProcessors(s map[string]map[string]string) *Processors {
	return &Processors{
		storeMap: s,
	}
}

// Get returns a data value from the store.
func (s *Processors) Get(processor string, key string) (string, bool) {
	col, ok := s.storeMap[processor]
	if !ok {
		return "", false
	}

	val, ok := col[key]
	return val, ok
}

// GetKeys returns all the keys for a specific processor from the store.
func (s *Processors) GetKeys(processor string) ([]string, bool) {
	col, ok := s.storeMap[processor]
	if !ok {
		return nil, false
	}

	keys := make([]string, 0, len(col))
	for k := range col {
		keys = append(keys, k)
	}

	return keys, true
}

// Set sets a config value in the store.
func (s *Processors) Set(processor string, key string, val string) error {
	col, ok := s.storeMap[processor]
	if !ok {
		s.storeMap[processor] = make(map[string]string)
		col = s.storeMap[processor]
	}

	col[key] = val
	return nil
}

// Delete deletes data from the store.
func (s *Processors) Delete(processor string, key string) {
	col, ok := s.storeMap[processor]
	if !ok {
		return
	}

	delete(col, key)
}

// DeleteAll deletes all data for a specific processor from the store.
func (s *Processors) DeleteAll(processor string) {
	col, ok := s.storeMap[processor]
	if !ok {
		return
	}

	for key := range col {
		delete(col, key)
	}

	delete(s.storeMap, processor)
}

// Exists returns true if there is any data for a specific processor in the store.
func (s *Processors) Exists(processor string) bool {
	_, ok := s.storeMap[processor]
	return ok
}
