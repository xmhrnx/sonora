// Package decode reads audio files - WAV, MP3, FLAC, Ogg Opus and Ogg
// Vorbis natively, anything else through ffmpeg where it is installed - as
// mono samples at a chosen rate. It decodes a block at a time and
// resamples as it goes, so an hour-long recording costs memory only at
// the target rate.
package decode

import (
	"bufio"
	"bytes"
	"errors"
	"os"

	"github.com/mehran9675/sonora/internal/audio"
	"github.com/mehran9675/sonora/internal/audio/opusfile"
	"github.com/mehran9675/sonora/internal/audio/wav"
)

// SpeechRate is the sample rate speech recognizers take.
const SpeechRate = 16000

// Extensions are the audio files worth offering in a file picker.
var Extensions = []string{"wav", "mp3", "flac", "ogg", "oga", "opus", "m4a", "aac", "webm", "mp4", "aiff", "wma"}

// ErrFormat means the file is in a format only ffmpeg can read and
// ffmpeg isn't installed (or, on a phone, can't be).
var ErrFormat = errors.New("this audio format can't be read here; use WAV, MP3, FLAC, Opus or Ogg Vorbis, or install ffmpeg")

// File decodes the audio file at path as mono samples at rate, or at the
// file's own rate when rate is 0.
func File(path string, rate int) (*audio.Clip, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 64<<10)
	head, _ := r.Peek(64)
	out := newResampler(rate)
	switch {
	case wav.IsWAV(head):
		err = wav.Stream(r, out.write)
	case bytes.HasPrefix(head, []byte("fLaC")):
		err = decodeFLAC(r, out)
	case opusfile.IsOpus(head):
		err = opusfile.Decode(r, out.write)
	case bytes.HasPrefix(head, []byte("OggS")):
		err = decodeVorbis(r, out)
	case bytes.HasPrefix(head, []byte("ID3")), len(head) >= 2 && head[0] == 0xFF && head[1]&0xE0 == 0xE0:
		err = decodeMP3(r, out) // fails on ADTS AAC, which looks alike
	default:
		return viaFFmpeg(path, rate)
	}
	if err != nil && out.empty() {
		if clip, ferr := viaFFmpeg(path, rate); ferr == nil {
			return clip, nil
		}
	}
	if err != nil {
		return nil, err
	}
	clip := out.clip()
	if len(clip.Samples) == 0 {
		return nil, errors.New("the file has no audio")
	}
	return clip, nil
}
