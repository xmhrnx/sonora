//go:build !windows && !ios

package record

// #include <stdlib.h>
import "C"
import (
	"errors"
	"time"

	"github.com/gen2brain/malgo"
)

// miniaudio captures through CoreAudio on macOS, PulseAudio, ALSA or
// JACK on Linux, and AAudio or OpenSL ES on Android (the app must hold
// the RECORD_AUDIO permission first).
type miniaudio struct {
	ctx *malgo.AllocatedContext
	dev *malgo.Device
}

// firstData is how long a backend may take to deliver audio before the
// next one is tried.
const firstData = 800 * time.Millisecond

// errNotConnected means the chosen microphone isn't there.
var errNotConnected = errors.New("the microphone isn't connected")
