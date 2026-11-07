//go:build !windows && !ios

package record

// #include <stdlib.h>
import "C"
import (
	"runtime"
	"strings"

	"github.com/gen2brain/malgo"
)

// backends lists what to try, in order. A backend can open fine and
// still deliver nothing (seen with some PulseAudio setups), so each is
// checked for actual audio.
func backends() [][]malgo.Backend {
	if runtime.GOOS == "linux" {
		return [][]malgo.Backend{{malgo.BackendPulseaudio}, {malgo.BackendAlsa}, {malgo.BackendJack}}
	}
	return [][]malgo.Backend{nil} // the platform default
}

// Devices lists the microphones from the first backend that has any.
func Devices() ([]Device, error) {
	var lastErr error
	for _, b := range backends() {
		ctx, err := malgo.InitContext(b, malgo.ContextConfig{}, nil)
		if err != nil {
			lastErr = err
			continue
		}
		infos, err := ctx.Devices(malgo.Capture)
		ctx.Uninit()
		ctx.Free()
		if err != nil {
			lastErr = err
			continue
		}
		var out []Device
		seen := map[string]bool{}
		for i := range infos {
			name := infos[i].Name()
			// PulseAudio lists every output's loopback as an input too,
			// and ALSA its plugins (rate converters, JACK, PulseAudio…).
			if name == "" || seen[name] || strings.HasPrefix(name, "Monitor of ") ||
				isALSA(b) && !isCard(infos[i].ID) {
				continue
			}
			seen[name] = true
			out = append(out, Device{ID: name, Name: name})
		}
		if len(out) > 0 {
			return out, nil
		}
	}
	return nil, lastErr
}

func isALSA(b []malgo.Backend) bool { return len(b) == 1 && b[0] == malgo.BackendAlsa }

// isCard reports whether an ALSA device ID names a sound card
// ("hw:0,0", "sysdefault:CARD=PCH") rather than a plugin ("jack").
func isCard(id malgo.DeviceID) bool {
	name, _, _ := strings.Cut(string(id[:]), "\x00")
	return strings.Contains(name, ":")
}
