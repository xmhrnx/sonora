package playback

import (
	"github.com/mehran9675/sonora/internal/audio"
)

// Load replaces the current clip; key identifies it (usually the file
// path) so callers can tell what is loaded.
func (pl *Player) Load(c *audio.Clip, key string) error {
	dev, err := device()
	if err != nil {
		return err
	}
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.closeLocked()
	pl.src = &source{samples: c.Resample(deviceRate)}
	pl.p = dev.NewPlayer(pl.src)
	pl.p.SetBufferSize(deviceRate * frameBytes / 10)
	pl.length, pl.key = c.Duration(), key
	return nil
}

// Loaded returns the key of the loaded clip.
func (pl *Player) Loaded() string {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.key
}

// Unload stops and forgets the current clip.
func (pl *Player) Unload() {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	pl.closeLocked()
	pl.p, pl.src, pl.key, pl.length = nil, nil, "", 0
}

func (pl *Player) closeLocked() {
	if pl.p != nil {
		pl.p.Pause()
		pl.p.Close()
	}
}
