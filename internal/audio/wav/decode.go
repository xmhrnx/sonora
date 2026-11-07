// Package wav reads and writes WAV files.
package wav

import (
	"bytes"
	"errors"
	"os"

	"github.com/mehran9675/sonora/internal/audio"
)

type format struct {
	tag, channels, bits uint16
	rate                uint32
}

// Read decodes a WAV file (PCM 8/16/24/32-bit or IEEE float), mixing
// every channel down to mono.
func Read(path string) (*audio.Clip, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Decode(data)
}

// IsWAV reports whether data starts with a RIFF/WAVE header.
func IsWAV(data []byte) bool {
	return len(data) >= 12 && string(data[0:4]) == "RIFF" && string(data[8:12]) == "WAVE"
}

// Decode decodes WAV bytes; see Read.
func Decode(data []byte) (*audio.Clip, error) {
	if !IsWAV(data) {
		return nil, errors.New("not a WAV file")
	}
	f, pcm, err := chunks(bytes.NewReader(data[12:]))
	if err != nil {
		return nil, err
	}
	bps := int(f.bits) / 8
	frame := bps * int(f.channels)
	if frame == 0 || f.rate == 0 {
		return nil, errors.New("bad WAV header")
	}
	n := len(pcm) / frame
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		var sum float32
		for c := 0; c < int(f.channels); c++ {
			off := i*frame + c*bps
			v, err := sample(f, pcm[off:off+bps])
			if err != nil {
				return nil, err
			}
			sum += v
		}
		out[i] = sum / float32(f.channels)
	}
	return &audio.Clip{Samples: out, SampleRate: int(f.rate)}, nil
}
