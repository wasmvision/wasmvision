package runtime

import (
	"math/rand/v2"
	"sync"
)

const maxIndex = 1048560

// MapRefs is a [Refs] implementation powered by a map protected by a mutex.
// Indexes are generated randomly and checked for collisions.
//
// Must be constructed with [NewMapRefs].
type MapRefs struct {
	mux sync.Mutex

	Raw map[uint32]any
	idx uint32
}

// NewMapRefs creates a new MapRefs.
func NewMapRefs() *MapRefs {
	return &MapRefs{Raw: make(map[uint32]any)}
}

// Get returns a value by index.
func (r *MapRefs) Get(idx uint32, def any) (any, bool) {
	r.mux.Lock()
	defer r.mux.Unlock()

	val, found := r.Raw[idx]
	if !found {
		return def, false
	}
	return val, true
}

// Set sets a value by index.
func (r *MapRefs) Set(idx uint32, val any) {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.Raw[idx] = val
}

// Put puts a value into the map and returns its index.
func (r *MapRefs) Put(val any) uint32 {
	r.mux.Lock()
	defer r.mux.Unlock()

	r.idx = uint32(rand.IntN(maxIndex))

	// skip already used cells
	_, used := r.Raw[r.idx]
	for used {
		r.idx = uint32(rand.IntN(maxIndex))
		_, used = r.Raw[r.idx]
	}

	r.Raw[r.idx] = val
	return r.idx
}

// Drop removes a value by index.
func (r *MapRefs) Drop(idx uint32) {
	r.mux.Lock()
	defer r.mux.Unlock()

	delete(r.Raw, idx)
}
