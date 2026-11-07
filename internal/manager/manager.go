// Package manager connects engines to the UI: it builds the engine
// list, runs installs, voice loading, generation and transcription in
// the background, and publishes their progress as Runtime values on the
// UI goroutine.
package manager

import (
	"github.com/mehran9675/sonora/internal/library"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/internal/transcripts"
)

// Manager owns every engine's Runtime.
type Manager struct {
	post        func(func())
	settings    *settings.Store
	library     *library.Library
	transcripts *transcripts.Store
	script      string

	list []*Runtime
	byID map[string]*Runtime

	// Job is the generation in progress, if any; TranscribeJob the
	// transcription. One of each can run at a time.
	Job           *Job
	TranscribeJob *Job
}
