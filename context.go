package dispatcher

import (
	"sync"
	"time"
)

var (
	mutex sync.RWMutex
	data  = make(map[*State]interface{})
	datat = make(map[*State]int64)
)

// Set stores a context for particular state.
func Set(s *State, val interface{}) {
	mutex.Lock()
	if data[s] == nil {
		datat[s] = time.Now().Unix()
	}
	data[s] = val
	mutex.Unlock()
}

// Get returns the context stores for particular state.
func Get(s *State) interface{} {
	mutex.RLock()
	if ctx := data[s]; ctx != nil {
		value := ctx
		mutex.RUnlock()
		return value
	}
	mutex.RUnlock()
	return nil
}
