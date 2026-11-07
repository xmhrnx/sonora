package playback

import (
	"io"
	"time"
)

// Play starts or resumes playback, rewinding a finished clip.
func (pl *Player) Play() {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if pl.p == nil {
		return
	}
	if pl.src.position() >= pl.src.total() && pl.p.BufferedSize() == 0 {
		pl.p.Seek(0, io.SeekStart)
	}
	pl.p.Play()
}

// Pause pauses playback.
func (pl *Player) Pause() {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if pl.p != nil {
		pl.p.Pause()
	}
}

// Toggle plays or pauses.
func (pl *Player) Toggle() {
	if pl.Playing() {
		pl.Pause()
	} else {
		pl.Play()
	}
}

// Playing reports whether audio is being played.
func (pl *Player) Playing() bool {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.p != nil && pl.p.IsPlaying()
}

// Length returns the loaded clip's duration.
func (pl *Player) Length() time.Duration {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	return pl.length
}

// Progress returns the playback position as a fraction in [0, 1].
func (pl *Player) Progress() float32 {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if pl.p == nil || pl.src.total() == 0 {
		return 0
	}
	pos := pl.src.position() - int64(pl.p.BufferedSize())
	return float32(max(0, min(1, float64(pos)/float64(pl.src.total()))))
}

// Seek moves playback to the fraction f of the clip.
func (pl *Player) Seek(f float32) {
	pl.mu.Lock()
	defer pl.mu.Unlock()
	if pl.p == nil {
		return
	}
	pl.p.Seek(int64(float64(pl.src.total())*float64(max(0, min(1, f)))), io.SeekStart)
}
