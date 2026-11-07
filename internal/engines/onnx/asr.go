package onnx

import (
	"context"
	"errors"
	"sync"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// Recognizer is an on-device speech-to-text engine.
type Recognizer struct {
	m      *engines.Manifest
	dir    string
	pick   func(engines.Values) model
	models []model

	mu     sync.Mutex
	rec    *sherpa.Recognizer
	vad    *sherpa.VAD
	loaded string // what rec was created for, see key
}

// maxUtterance is the longest stretch, in seconds, handed to a
// recognizer at once; Whisper takes at most 30.
const maxUtterance = 20

func newRecognizer(m *engines.Manifest, dir string, models []model) *Recognizer {
	pick := byOption(models)
	return &Recognizer{m: m, dir: dir, models: models, pick: func(o engines.Values) model { return pick("", o) }}
}

func (e *Recognizer) Manifest() *engines.Manifest { return e.m }

// Models lists the packages to choose from when the manifest has no
// model option (an engine with a single model).
func (e *Recognizer) Models() []engines.Choice { return choices(e.models) }

// Installed reports whether any of the engine's models is downloaded.
func (e *Recognizer) Installed() bool { return anyDownloaded(e.dir, e.models) }

// DownloadedModels reports which packages are on this device.
func (e *Recognizer) DownloadedModels() map[string]bool { return downloaded(e.dir, e.models) }

// InstallModel downloads the package with the given id.
func (e *Recognizer) InstallModel(ctx context.Context, id string, log engines.Logger) error {
	return fetchModel(ctx, e.dir, e.models, id, log)
}

// Install downloads nothing; see Engine.Install.
func (e *Recognizer) Install(context.Context, engines.Logger) error { return errChooseModel }

// Uninstall unloads and deletes every package this engine downloaded.
func (e *Recognizer) Uninstall() error {
	e.Close()
	return removeModels(e.dir, e.models)
}

func (e *Recognizer) Voices(context.Context, engines.Values, engines.Logger) ([]engines.Voice, error) {
	return nil, nil
}

func (e *Recognizer) Synthesize(context.Context, engines.Request, engines.Logger) error {
	return errors.New(e.m.Name + " turns speech into text; it can't speak")
}

// Loaded reports whether a model is in memory.
func (e *Recognizer) Loaded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.rec != nil
}

// Close frees the model's memory.
func (e *Recognizer) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.unload()
}

func (e *Recognizer) unload() {
	if e.rec != nil {
		e.rec.Close()
		e.rec, e.loaded = nil, ""
	}
	if e.vad != nil {
		e.vad.Close()
		e.vad = nil
	}
}
