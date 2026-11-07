// Package record captures audio from a microphone - a voice
// to clone, or speech to transcribe - straight into an Ogg Opus file,
// so a recording's length is bounded by the disk, not memory.
package record

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/mehran9675/sonora/internal/audio/opusfile"
)

// capture is the platform's input device.
type capture interface {
	stop()
}

// Options say where and how to record.
type Options struct {
	Path    string // the .opus file to write; its folder is created
	Bitrate int    // Opus bitrate in bits per second
	// Device is the microphone's Device.ID; empty, or a microphone
	// that isn't connected, records from the system default.
	Device string
}

// Device is a microphone.
type Device struct {
	// ID is what Options.Device takes. Desktops know a microphone by
	// its name, which outlives restarts and backends where their IDs
	// don't.
	ID   string
	Name string
}

// Recorder records mono audio into an Opus file. Methods are safe from
// any goroutine.
type Recorder struct {
	path string
	w    *opusfile.Writer
	dev  capture
	done chan error // the writer goroutine's result

	mu       sync.Mutex
	wake     *sync.Cond
	queue    [][]float32 // captured blocks waiting to be encoded
	closed   bool
	captured int64 // samples captured
	levels   []float32
}

// SampleRate is the capture rate, Opus's own.
const SampleRate = 48000

// levelWindow is how many recent loudness values Levels keeps.
const levelWindow = 48

// ErrNoInput means no microphone could be opened.
var ErrNoInput = errors.New("no microphone found, or access to it was denied")

// AndroidPrefix starts the ID of an Android input device: the prefix
// and AudioDeviceInfo.getId(). miniaudio lists only the default input
// on Android, so the app finds the others through AudioManager.
const AndroidPrefix = "android:"

// active counts recordings capturing from a microphone.
var active atomic.Int32

// Active reports whether a recording is running.
func Active() bool { return active.Load() > 0 }
