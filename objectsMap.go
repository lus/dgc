package dgc

import "sync"

// ObjectsMap wraps a map[string]interface
// to provide thread safe access endpoints.
type ObjectsMap struct {
	mtx sync.RWMutex
	m   map[string]interface{}
}

// newObjectsMap initializes a new
// ObjectsMap instance
func newObjectsMap() *ObjectsMap {
	return &ObjectsMap{
		m: make(map[string]interface{}),
	}
}

// Get tries to get a value from the map.
// If a value was found, the value and true
// is returned. Else, nil and false is
// returned.
func (om *ObjectsMap) Get(key string) (interface{}, bool) {
	om.mtx.RLock()
	defer om.mtx.RUnlock()

	v, ok := om.m[key]
	return v, ok
}

// MustGet wraps Get but only returns the
// value, if found, or nil otherwise.
func (om *ObjectsMap) MustGet(key string) interface{} {
	v, ok := om.Get(key)
	if !ok {
		return nil
	}
	return v
}

// Set sets a value to the map by key.
func (om *ObjectsMap) Set(key string, val interface{}) {
	om.mtx.Lock()
	defer om.mtx.Unlock()

	om.m[key] = val
}

// Delete removes a key-value pair from the map.
func (om *ObjectsMap) Delete(key string) {
	om.mtx.Lock()
	defer om.mtx.Unlock()

	delete(om.m, key)
}
