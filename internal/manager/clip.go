package manager

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/audio/wav"
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/paths"
)

// TranscribeClip hears what is said in a short recording without
// keeping a transcript - a voice profile made from an imported file.
// The language is detected. done runs on the UI goroutine.
func (m *Manager) TranscribeClip(rt *Runtime, path string, done func(text string, err error)) {
	tr, ok := rt.Engine.(engines.Transcriber)
	man := rt.Manifest()
	switch {
	case !ok || !man.Transcribes():
		done("", errors.New(man.Name+" can't transcribe"))
		return
	case m.TranscribeJob != nil || rt.Busy():
		done("", errors.New("already transcribing"))
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	rt.cancel = cancel
	rt.State, rt.Err, rt.Progress, rt.Status = StateWorking, "", -1, "Reading the audio…"
	m.TranscribeJob = &Job{Runtime: rt, Started: time.Now()}
	sink := m.sink(rt)
	go func() {
		text, err := m.transcribeClip(ctx, tr, man, path, sink.Log)
		m.post(func() {
			rt.cancel, m.TranscribeJob = nil, nil
			rt.State, rt.Progress = StateReady, -1
			switch {
			case errors.Is(err, context.Canceled):
				rt.Status = "Stopped"
			case err != nil:
				rt.Err, rt.Status = err.Error(), "Transcription failed"
			default:
				rt.Status = "Done"
			}
			done(text, err)
		})
	}()
}

func (m *Manager) transcribeClip(ctx context.Context, tr engines.Transcriber, man *engines.Manifest, path string,
	log engines.Logger) (string, error) {
	log("Reading the audio…")
	clip, err := decode.File(path, decode.SpeechRate)
	if err != nil {
		return "", err
	}
	input := filepath.Join(paths.Sub("cache"), "clip-"+strconv.FormatInt(time.Now().UnixNano(), 36)+".wav")
	if err := wav.Write(input, clip); err != nil {
		return "", err
	}
	defer os.Remove(input)
	res, err := tr.Transcribe(ctx, engines.TranscribeRequest{Audio: input, Options: m.values(man)}, log)
	if err != nil {
		return "", explain(err)
	}
	text := strings.TrimSpace(res.Text)
	if text == "" {
		return "", errors.New("no speech was recognized in the recording")
	}
	return text, nil
}
