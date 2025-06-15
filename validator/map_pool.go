package validator

import "sync"

// MapPool provides a pool of reusable string maps
var MapPool = &sync.Pool{
	New: func() interface{} {
		return make(map[string]string)
	},
}

// GetMap retrieves a map from the pool or creates a new one if none is available
func getMap() map[string]string {
	return MapPool.Get().(map[string]string)
}

// PutMap returns a map to the pool after clearing its contents
func PutMap(m map[string]string) {
	// Clear the map before returning it to the pool
	for k := range m {
		delete(m, k)
	}
	MapPool.Put(m)
}
