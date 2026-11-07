// Package transcripts keeps every transcription with its audio, as Opus.
package transcripts

import (
	"path/filepath"
	"sync"
	"time"

	"github.com/mehran9675/sonora/internal/engines"
)

// Entry is one transcribed recording.
type Entry struct {
	ID         string            `json:"id"`
	Title      string            `json:"title"` // the source file's name
	Text       string            `json:"text"`
	Language   string            `json:"language,omitempty"`
	Translated bool              `json:"translated,omitempty"`
	Segments   []engines.Segment `json:"segments,omitempty"`
	EngineID   string            `json:"engine_id"`
	EngineName string            `json:"engine_name"`
	Created    time.Time         `json:"created"`
	Duration   time.Duration     `json:"duration"`
	Audio      string            `json:"audio"` // the audio inside the folder: Opus (older entries: WAV)
}

// Store is the folder of transcripts plus index.json. Methods are safe
// from any goroutine.
type Store struct {
	mu      sync.Mutex
	dir     string
	entries []Entry
}

func (s *Store) index() string { return filepath.Join(s.dir, "index.json") }
