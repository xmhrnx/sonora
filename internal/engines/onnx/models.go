// Package onnx runs speech models inside the app with sherpa-onnx, on
// the CPU and with no Python: text to speech with Kitten, Kokoro, Piper
// and ZipVoice, speech to text with Whisper, Parakeet and SenseVoice.
// It's what Android uses; desktop builds include it when built with the
// "sherpa" tag.
package onnx

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// releases is where sherpa-onnx publishes its ready-to-run models.
const releases = "https://github.com/k2-fsa/sherpa-onnx/releases/download/"

const (
	// vocos is the vocoder ZipVoice models need, a separate download.
	vocos = releases + "vocoder-models/vocos_24khz.onnx"
	// sileroVAD finds the speech in recordings before recognition.
	sileroVAD = releases + "asr-models/silero_vad.onnx"
)

// model is one downloadable package.
type model struct {
	ID     string // package name; the archive is <release>/<ID>.tar.bz2
	Label  string
	Size   string
	Family sherpa.Family
	// Release is the release holding the package: "tts-models" when
	// empty.
	Release string
	// Lexicon lists lexicon files inside the package (Kokoro).
	Lexicon string
	// Dict is the dictionary folder inside the package (Kokoro).
	Dict string
	// Steps is the number of flow-matching steps (ZipVoice).
	Steps int
	// Int8 unpacks only the int8 models of a package that ships both.
	Int8 bool
	// Extra are single files the model also needs, downloaded once
	// next to the packages and shared.
	Extra []string
}

func (m model) url() string {
	rel := m.Release
	if rel == "" {
		rel = "tts-models"
	}
	return releases + rel + "/" + m.ID + ".tar.bz2"
}

// keep reports whether a file of the package is unpacked.
func (m model) keep(name string) bool {
	if strings.HasPrefix(name, "test_wavs/") {
		return false
	}
	if m.Int8 && strings.HasSuffix(name, ".onnx") && !strings.HasSuffix(name, ".int8.onnx") {
		return false
	}
	return true
}

// modelOption lets the user pick one of an engine's packages.
func modelOption(models []model, help string) engines.Option {
	return engines.Option{ID: "model", Label: "Model", Type: engines.OptionSelect, Scope: engines.ScopeLoad,
		Default: models[0].ID, Choices: choices(models), Help: help}
}

// errChooseModel stops an install that didn't name a model.
var errChooseModel = errors.New("choose a model to download")

// choices lists packages for the user to pick from, with their sizes.
func choices(models []model) []engines.Choice {
	var out []engines.Choice
	for _, m := range models {
		out = append(out, engines.Choice{Value: m.ID, Label: m.Label + " · " + m.Size})
	}
	return out
}

// downloaded reports which of models are unpacked in dir.
func downloaded(dir string, models []model) map[string]bool {
	out := map[string]bool{}
	for _, m := range models {
		if _, err := os.Stat(filepath.Join(dir, m.ID)); err == nil {
			out[m.ID] = true
		}
	}
	return out
}

func anyDownloaded(dir string, models []model) bool { return len(downloaded(dir, models)) > 0 }

// fetchModel downloads the package of models with the given id.
func fetchModel(ctx context.Context, dir string, models []model, id string, log engines.Logger) error {
	for _, m := range models {
		if m.ID == id {
			_, err := fetch(ctx, dir, m, log)
			return err
		}
	}
	return fmt.Errorf("unknown model %q", id)
}

// removeModels deletes downloaded packages and their extra files.
func removeModels(dir string, models []model) error {
	var paths []string
	for _, m := range models {
		paths = append(paths, filepath.Join(dir, m.ID))
		for _, url := range m.Extra {
			paths = append(paths, extraPath(dir, url))
		}
	}
	for _, p := range paths {
		if err := os.RemoveAll(p); err != nil {
			return err
		}
	}
	return nil
}
