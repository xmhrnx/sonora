//go:build sherpa || android

package sherpa

/*
#include <stdlib.h>
#include "c-api.h"
*/
import "C"
import (
	"errors"
	"sync"
)

// TTS is a loaded model. Generate may be called from any goroutine,
// one call at a time.
type TTS struct {
	mu  sync.Mutex
	ptr *C.SherpaOnnxOfflineTts
}

// Available reports whether this build can run on-device models.
const Available = true

// ErrStopped is returned when progress asked to stop.
var ErrStopped = errors.New("generation stopped")
