package openai

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/mehran9675/sonora/internal/engines"
)

// Transcriber is one configured transcription endpoint.
type Transcriber struct {
	*Engine
}

// TranscribeTemplate returns the manifest for a new transcription API
// engine.
func TranscribeTemplate(id, name string) *engines.Manifest {
	return &engines.Manifest{
		ID:          id,
		Name:        name,
		Vendor:      "OpenAI-compatible API",
		Description: "Transcribes through an OpenAI-compatible /v1/audio/transcriptions endpoint - OpenAI, Groq, a local Whisper server…",
		Tags:        []string{"API", "Remote"},
		Accent:      []string{"#94a3b8", "#475569"},
		Kind:        engines.KindOpenAITranscribe,
		Task:        engines.TaskTranscribe,
		Translate:   true,
		Custom:      true,
		Options: []engines.Option{
			{ID: "base_url", Label: "Server URL", Type: engines.OptionText, Scope: engines.ScopeSetup,
				Default: "https://api.openai.com/v1", Placeholder: "https://api.openai.com/v1",
				Help: "Everything before /audio/transcriptions", Keyboard: "url"},
			{ID: "api_key", Label: "API key", Type: engines.OptionSecret, Scope: engines.ScopeSetup, Placeholder: "optional for local servers"},
			{ID: "model", Label: "Model", Type: engines.OptionText, Scope: engines.ScopeSetup, Default: "whisper-1",
				Suggestions: []string{"whisper-1", "gpt-4o-transcribe", "gpt-4o-mini-transcribe", "whisper-large-v3-turbo", "Systran/faster-whisper-large-v3"}},
			{ID: "prompt", Label: "Vocabulary hint", Type: engines.OptionText, Placeholder: "Names, jargon, spelling…"},
		},
	}
}

// NewTranscriber returns the engine; config returns its option values.
func NewTranscriber(m *engines.Manifest, config func() engines.Values) *Transcriber {
	return &Transcriber{&Engine{m: m, config: config, client: &http.Client{Timeout: 30 * time.Minute}}}
}

func (e *Transcriber) Voices(context.Context, engines.Values, engines.Logger) ([]engines.Voice, error) {
	return nil, nil
}

func (e *Transcriber) Synthesize(context.Context, engines.Request, engines.Logger) error {
	return errors.New(e.m.Name + " turns speech into text; it can't speak")
}

// maxUpload is OpenAI's file size limit.
const maxUpload = 25 << 20
