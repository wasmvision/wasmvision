package datastore

type Frames struct {
	storeMap map[int]map[string]string
}

// NewFrames creates a new Frames data store.
func NewFrames(s map[int]map[string]string) *Frames {
	return &Frames{
		storeMap: s,
	}
}

// Get returns a data value for a specific frame from the store .
func (s *Frames) Get(frame int, key string) (string, bool) {
	col, ok := s.storeMap[frame]
	if !ok {
		return "", false
	}

	val, ok := col[key]
	return val, ok
}

// GetKeys returns all the keys for a specific frame from the store.
func (s *Frames) GetKeys(frame int) ([]string, bool) {
	col, ok := s.storeMap[frame]
	if !ok {
		return nil, false
	}

	keys := make([]string, 0, len(col))
	for k := range col {
		keys = append(keys, k)
	}

	return keys, true
}

// Set sets a key/value for a specific frame in the store.
func (s *Frames) Set(frame int, key, val string) error {
	col, ok := s.storeMap[frame]
	if !ok {
		s.storeMap[frame] = make(map[string]string)
		col = s.storeMap[frame]
	}

	col[key] = val
	return nil
}

// Delete deletes data for a specific frame from the store.
func (s *Frames) Delete(frame int, key string) {
	col, ok := s.storeMap[frame]
	if !ok {
		return
	}

	delete(col, key)
}

// DeleteAll deletes all data for a specific frame from the store.
func (s *Frames) DeleteAll(frame int) {
	col, ok := s.storeMap[frame]
	if !ok {
		return
	}

	for key := range col {
		delete(col, key)
	}

	delete(s.storeMap, frame)
}

// Exists returns true if there is any data for a specific frame in the store.
func (s *Frames) Exists(frame int) bool {
	_, ok := s.storeMap[frame]
	return ok
}
