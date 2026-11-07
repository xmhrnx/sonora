package onnx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mehran9675/sonora/internal/audio"
	"github.com/mehran9675/sonora/internal/audio/convert"
	"github.com/mehran9675/sonora/internal/audio/wav"
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// load returns the model for m, downloading and loading it as needed.
func (e *Engine) load(ctx context.Context, m model, log engines.Logger) (*sherpa.TTS, error) {
	dir, err := fetch(ctx, e.dir, m, log)
	if err != nil {
		return nil, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.tts != nil && e.loaded == m.ID {
		return e.tts, nil
	}
	if e.tts != nil {
		e.tts.Close()
		e.tts = nil
	}
	cfg, err := config(dir, m)
	if err != nil {
		return nil, err
	}
	if m.Family == sherpa.ZipVoice {
		cfg.Vocoder = extraPath(e.dir, vocos)
	}
	log("Loading " + e.m.Name + "…")
	tts, err := sherpa.New(cfg)
	if err != nil {
		return nil, err
	}
	e.tts, e.loaded = tts, m.ID
	return tts, nil
}

// config finds a package's files.
func config(dir string, m model) (sherpa.Config, error) {
	find := func(pattern string) (string, error) {
		found, _ := filepath.Glob(filepath.Join(dir, pattern))
		if len(found) == 0 {
			return "", errors.New("the downloaded model is incomplete; remove the engine and install it again")
		}
		return found[0], nil
	}
	cfg := sherpa.Config{Family: m.Family,
		Tokens: filepath.Join(dir, "tokens.txt"), DataDir: filepath.Join(dir, "espeak-ng-data")}
	var err error
	if m.Family == sherpa.ZipVoice {
		if cfg.Encoder, err = find("encoder*.onnx"); err != nil {
			return cfg, err
		}
		if cfg.Decoder, err = find("decoder*.onnx"); err != nil {
			return cfg, err
		}
	} else if cfg.Model, err = find("*.onnx"); err != nil {
		return cfg, err
	}
	if m.Family == sherpa.Kitten || m.Family == sherpa.Kokoro {
		cfg.Voices = filepath.Join(dir, "voices.bin")
	}
	if m.Lexicon != "" {
		cfg.Lexicon = filepath.Join(dir, m.Lexicon)
	}
	if m.Dict != "" {
		cfg.DictDir = filepath.Join(dir, m.Dict)
	}
	return cfg, nil
}

// Synthesize generates speech and writes it as WAV.
func (e *Engine) Synthesize(ctx context.Context, r engines.Request, log engines.Logger) error {
	m := e.pick(r.Voice, r.Options)
	p := sherpa.Params{Speed: float32(r.Speed), Steps: m.Steps}
	if p.Speed <= 0 {
		p.Speed = 1
	}
	if m.Family == sherpa.Kitten || m.Family == sherpa.Kokoro {
		p.Speaker, _ = strconv.Atoi(r.Voice)
	}
	if m.Family == sherpa.ZipVoice {
		if r.Reference == "" || r.ReferenceText == "" {
			return errors.New(e.m.Name + " needs a recording of the voice and its transcript")
		}
		clip, err := readReference(r.Reference)
		if err != nil {
			return err
		}
		p.Reference, p.ReferenceRate, p.ReferenceText = clip.Samples, clip.SampleRate, r.ReferenceText
	}
	tts, err := e.load(ctx, m, log)
	if err != nil {
		return err
	}
	last := -1
	samples, rate, err := tts.Generate(r.Text, p, func(p float32) bool {
		if pct := int(p * 100); pct != last {
			last = pct
			log(fmt.Sprintf("Generating… %d%%", pct))
		}
		return ctx.Err() == nil
	})
	if errors.Is(err, sherpa.ErrStopped) && ctx.Err() != nil {
		return ctx.Err()
	}
	if err != nil {
		return err
	}
	return wav.Write(r.Output, &audio.Clip{Samples: samples, SampleRate: rate})
}

// readReference loads a cloning recording. Other formats than WAV are
// converted with ffmpeg where it is installed (not on phones; recordings
// made in the app are WAV).
func readReference(path string) (*audio.Clip, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if !wav.IsWAV(data) {
		if data, err = convert.ToWAV(data); err != nil {
			return nil, errors.New("this engine needs a WAV recording; record your voice in the app or pick a WAV file")
		}
	}
	return wav.Decode(data)
}
