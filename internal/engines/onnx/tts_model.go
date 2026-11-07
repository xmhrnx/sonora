package onnx

import (
	"context"

	"github.com/mehran9675/sonora/internal/engines"
)

// Models lists the packages to choose from when the manifest has no
// model option (Piper: one per voice).
func (e *Engine) Models() []engines.Choice { return choices(e.models) }

// Installed reports whether any of the engine's models is downloaded.
func (e *Engine) Installed() bool { return anyDownloaded(e.dir, e.models) }

// Install downloads nothing: on-device engines are their models, and
// the user picks which to download (InstallModel).
func (e *Engine) Install(context.Context, engines.Logger) error { return errChooseModel }

// DownloadedModels reports which packages are on this device.
func (e *Engine) DownloadedModels() map[string]bool { return downloaded(e.dir, e.models) }

// InstallModel downloads the package with the given id.
func (e *Engine) InstallModel(ctx context.Context, id string, log engines.Logger) error {
	return fetchModel(ctx, e.dir, e.models, id, log)
}

// Uninstall unloads and deletes every package this engine downloaded.
func (e *Engine) Uninstall() error {
	e.Close()
	return removeModels(e.dir, e.models)
}

// Voices returns the fixed voice list.
func (e *Engine) Voices(context.Context, engines.Values, engines.Logger) ([]engines.Voice, error) {
	return e.m.Voices, nil
}

// Loaded reports whether a model is in memory.
func (e *Engine) Loaded() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.tts != nil
}

// Close frees the model's memory.
func (e *Engine) Close() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.tts != nil {
		e.tts.Close()
		e.tts, e.loaded = nil, ""
	}
}
