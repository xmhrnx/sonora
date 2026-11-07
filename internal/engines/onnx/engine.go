package onnx

import (
	"path/filepath"
	"sync"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// Engine is one on-device engine: Kitten, Kokoro, Piper or ZipVoice.
type Engine struct {
	m   *engines.Manifest
	dir string // downloaded packages: dir/<package id>/
	// pick returns the package a request needs.
	pick func(voice string, opts engines.Values) model
	// models are all the packages the engine may download.
	models []model

	mu     sync.Mutex
	tts    *sherpa.TTS
	loaded string // package id of tts
}

// Available reports whether this build can run on-device models.
const Available = sherpa.Available

// All returns the on-device engines, storing models under dataDir.
func All(dataDir string) []engines.Engine {
	dir := filepath.Join(dataDir, "models", "ondevice")
	var piper []model
	for _, v := range piperVoices {
		piper = append(piper, piperModel(v.ID))
	}
	return []engines.Engine{
		&Engine{m: kittenManifest(), dir: dir, pick: byOption(kittenModels), models: kittenModels},
		&Engine{m: kokoroManifest(), dir: dir, pick: byOption(kokoroModels), models: kokoroModels},
		&Engine{m: piperManifest(), dir: dir, models: piper, pick: func(voice string, _ engines.Values) model {
			if voice == "" {
				voice = piperVoices[0].ID
			}
			return piperModel(voice)
		}},
		&Engine{m: zipvoiceManifest(), dir: dir, pick: byOption(zipvoiceModels), models: zipvoiceModels},
		newRecognizer(whisperManifest(), dir, whisperModels),
		newRecognizer(parakeetManifest(), dir, parakeetModels),
		newRecognizer(senseVoiceManifest(), dir, senseVoiceModels),
		newRecognizer(moonshineManifest(), dir, moonshineModels),
		newRecognizer(barkManifest(), dir, barkModels),
	}
}

// byOption picks the package named by the "model" option.
func byOption(models []model) func(string, engines.Values) model {
	return func(_ string, o engines.Values) model {
		for _, m := range models {
			if m.ID == o.String("model") {
				return m
			}
		}
		return models[0]
	}
}

func (e *Engine) Manifest() *engines.Manifest { return e.m }
