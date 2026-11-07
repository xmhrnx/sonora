package decode

import (
	"bufio"
	"os"
	"path/filepath"

	"github.com/mehran9675/sonora/internal/audio"
	"github.com/mehran9675/sonora/internal/audio/convert"
	"github.com/mehran9675/sonora/internal/audio/wav"
)

// viaFFmpeg converts the file to WAV with ffmpeg and decodes that.
func viaFFmpeg(path string, rate int) (*audio.Clip, error) {
	if !convert.Available() {
		return nil, ErrFormat
	}
	tmp, err := os.MkdirTemp("", "sonora-decode-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	out := filepath.Join(tmp, "audio.wav")
	if err := convert.FileToWAV(path, out); err != nil {
		return nil, err
	}
	f, err := os.Open(out)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	res := newResampler(rate)
	if err := wav.Stream(bufio.NewReader(f), res.write); err != nil {
		return nil, err
	}
	return res.clip(), nil
}
