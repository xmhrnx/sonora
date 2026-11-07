package openai

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/mehran9675/sonora/internal/engines"
)

// Voices asks the server (GET /audio/voices, offered by several local
// servers) and falls back to the configured list.
func (e *Engine) Voices(ctx context.Context, _ engines.Values, log engines.Logger) ([]engines.Voice, error) {
	if vs := e.serverVoices(ctx); len(vs) > 0 {
		return vs, nil
	}
	var out []engines.Voice
	for _, v := range strings.Split(e.config().String("voices"), ",") {
		if v = strings.TrimSpace(v); v != "" {
			out = append(out, engines.Voice{ID: v, Name: v})
		}
	}
	return out, nil
}

func (e *Engine) serverVoices(ctx context.Context) []engines.Voice {
	ctx, cancel := context.WithTimeout(ctx, 4*time.Second)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, e.url("/audio/voices"), nil)
	resp, err := e.do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	var body struct {
		Voices []json.RawMessage `json:"voices"`
	}
	if json.NewDecoder(resp.Body).Decode(&body) != nil {
		return nil
	}
	var out []engines.Voice
	for _, raw := range body.Voices {
		var name string
		var obj struct{ ID, Name string }
		if json.Unmarshal(raw, &name) == nil {
			out = append(out, engines.Voice{ID: name, Name: name})
		} else if json.Unmarshal(raw, &obj) == nil && (obj.ID != "" || obj.Name != "") {
			id := obj.ID
			if id == "" {
				id = obj.Name
			}
			out = append(out, engines.Voice{ID: id, Name: obj.Name})
		}
	}
	return out
}
