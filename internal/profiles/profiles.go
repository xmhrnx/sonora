// Package profiles keeps voice profiles: voices to clone, each a short
// recording of one person and what they say in it, so the engine gets
// both without the user writing anything.
package profiles

import (
	"time"
)

// Profile is a voice to clone.
type Profile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Audio is the recording inside the folder, Opus.
	Audio string `json:"audio"`
	// Transcript is what is said in the recording: the passage read
	// aloud, or what a transcription engine heard in an imported file.
	// Empty until known.
	Transcript string        `json:"transcript,omitempty"`
	Duration   time.Duration `json:"duration"`
	// Imported marks profiles made from a file rather than recorded.
	Imported bool      `json:"imported,omitempty"`
	Created  time.Time `json:"created"`
}

// MaxLength is the longest recording a profile keeps; engines clone
// from 5-30 seconds and longer imports are cut.
const MaxLength = 30 * time.Second

// bitrate keeps voices clear enough to clone.
const bitrate = 64000

// rate is the rate profiles are stored at.
const rate = 48000
