package onnx

import (
	"context"
	"fmt"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// load returns the recognizer and voice detector for a request,
// downloading and (re)loading as needed. Whisper's language and task
// are fixed when it loads, so changing them reloads it.
func (e *Recognizer) load(ctx context.Context, r engines.TranscribeRequest, log engines.Logger) (*sherpa.Recognizer, *sherpa.VAD, error) {
	m := e.pick(r.Options)
	dir, err := fetch(ctx, e.dir, m, log)
	if err != nil {
		return nil, nil, err
	}
	cfg, err := recognizerConfig(dir, m)
	if err != nil {
		return nil, nil, err
	}
	cfg.Language = supportedLanguage(m, r.Language)
	cfg.Translate = r.Translate && m.Family == sherpa.Whisper
	key := fmt.Sprintf("%s|%s|%v", m.ID, cfg.Language, cfg.Translate)
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.rec != nil && e.loaded == key {
		return e.rec, e.vad, nil
	}
	e.unload()
	log("Loading " + e.m.Name + "…")
	vad, err := sherpa.NewVAD(extraPath(e.dir, sileroVAD), maxUtterance)
	if err != nil {
		return nil, nil, err
	}
	rec, err := sherpa.NewRecognizer(cfg)
	if err != nil {
		vad.Close()
		return nil, nil, err
	}
	e.rec, e.vad, e.loaded = rec, vad, key
	return rec, vad, nil
}
