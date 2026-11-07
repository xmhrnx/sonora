package profiles

import (
	"sync"
)

// Store is the folder of profiles plus index.json. Methods are safe
// from any goroutine.
type Store struct {
	mu   sync.Mutex
	dir  string
	list []Profile
}
