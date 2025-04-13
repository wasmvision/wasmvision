package storage

// MemStorage is an in-memory data storage implementation of the DataStorage interface.
// It uses a map to store data, where the key is a string and the value is another map of string keys and values.
type MemStorage[T int | string] struct {
	storeMap map[T]map[string]string
}

// NewMemStorage creates a new MemStorage data store.
func NewMemStorage[T int | string]() *MemStorage[T] {
	return &MemStorage[T]{
		storeMap: make(map[T]map[string]string),
	}
}

// Get returns a data value for a specific key from the store.
func (s *MemStorage[T]) Get(key T, subkey string) (string, bool) {
	col, ok := s.storeMap[key]
	if !ok {
		return "", false
	}

	val, ok := col[subkey]
	return val, ok
}

// GetKeys returns all the keys for a specific key from the store.
func (s *MemStorage[T]) GetKeys(key T) ([]string, bool) {
	col, ok := s.storeMap[key]
	if !ok {
		return nil, false
	}

	keys := make([]string, 0, len(col))
	for k := range col {
		keys = append(keys, k)
	}

	return keys, true
}

// Set sets a key/value for a specific key in the store.
func (s *MemStorage[T]) Set(key T, subkey, val string) error {
	col, ok := s.storeMap[key]
	if !ok {
		s.storeMap[key] = make(map[string]string)
		col = s.storeMap[key]
	}

	col[subkey] = val
	return nil
}

// Delete deletes data for a specific key from the store.
func (s *MemStorage[T]) Delete(key T, subkey string) {
	col, ok := s.storeMap[key]
	if !ok {
		return
	}

	delete(col, subkey)
}

// DeleteAll deletes all data for a specific key from the store.
func (s *MemStorage[T]) DeleteAll(key T) {
	delete(s.storeMap, key)
}

// Exists checks if a specific key exists in the store.
func (s *MemStorage[T]) Exists(key T) bool {
	_, ok := s.storeMap[key]
	return ok
}
