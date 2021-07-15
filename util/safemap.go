package util

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	bm   map[string]interface{}
}

// NewBeeMap return new safemap
func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[string]interface{}),
	}
}

func (m *SafeMap) Length() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.bm)
}

// Get from maps return the k's value
func (m *SafeMap) Get(k string) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.bm[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns false
// if the key is already in the map and changes nothing.
func (m *SafeMap) Set(k string, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.bm[k]; !ok {
		m.bm[k] = v
	} else if val != v {
		m.bm[k] = v
	} else {
		return false
	}
	return true
}

// Returns true if k is exist in the map.
func (m *SafeMap) Check(k string) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if _, ok := m.bm[k]; !ok {
		return false
	}
	return true
}

// Delete the given key and value.
func (m *SafeMap) Delete(k string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, k)
}

// Items returns all items in safemap.
func (m *SafeMap) Items() map[string]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[string]interface{})
	for k, v := range m.bm {
		r[k] = v
	}
	return r
}
// Clear clears all items in safemap.
// func (m *SafeMap) Clear() {
// 	m.lock.RLock()
// 	defer m.lock.RUnlock()
// 	for k, _ := range m.bm {
// 		delete(m.bm, k)
// 	}
// }
