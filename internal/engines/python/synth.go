package python

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mehran9675/sonora/internal/engines"
)

// Voices returns the manifest's voices, or loads the model and asks
// the adapter for them.
func (e *Engine) Voices(ctx context.Context, opts engines.Values, log engines.Logger) ([]engines.Voice, error) {
	if len(e.m.Voices) > 0 {
		return e.m.Voices, nil
	}
	if _, err := e.ensure(ctx, opts, log); err != nil {
		return nil, err
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.voices, nil
}

// Synthesize runs one request through the adapter.
func (e *Engine) Synthesize(ctx context.Context, req engines.Request, log engines.Logger) error {
	w, err := e.ensure(ctx, engines.Subset(e.m.Options, req.Options, engines.ScopeLoad), log)
	if err != nil {
		return err
	}
	_, err = w.call(ctx, "synthesize", nil, req)
	return err
}

// Transcribe runs one speech-to-text request through the adapter.
func (e *Engine) Transcribe(ctx context.Context, req engines.TranscribeRequest, log engines.Logger) (*engines.Transcript, error) {
	w, err := e.ensure(ctx, engines.Subset(e.m.Options, req.Options, engines.ScopeLoad), log)
	if err != nil {
		return nil, err
	}
	raw, err := w.call(ctx, "transcribe", nil, req)
	if err != nil {
		return nil, err
	}
	var t engines.Transcript
	if err := json.Unmarshal(raw, &t); err != nil {
		return nil, errors.New("the engine returned an unreadable transcript")
	}
	return &t, nil
}

// Loaded reports whether a worker with a loaded model is running.
func (e *Engine) Loaded() bool {
	w := e.current.Load()
	return w != nil && w.alive()
}

// Close stops the worker, freeing the model's memory.
func (e *Engine) Close() {
	if w := e.current.Swap(nil); w != nil {
		// Stop it first, so a Close during a long load or generation
		// doesn't wait for the lock that work holds.
		w.close()
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.w != nil {
		e.w.close()
		e.w = nil
	}
}
