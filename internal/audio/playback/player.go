package playback

import (
	"sync"
	"time"

	"github.com/ebitengine/oto/v3"
)

// Player plays one clip at a time. All methods are safe from any
// goroutine.
type Player struct {
	mu     sync.Mutex
	src    *source
	p      *oto.Player
	length time.Duration
	key    string
}

// New returns an idle player.
func New() *Player { return &Player{} }
