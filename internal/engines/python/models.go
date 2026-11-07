package python

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/mehran9675/sonora/internal/engines"
)

func (e *Engine) modelsPath() string { return filepath.Join(e.dir, "models.json") }

// DownloadedModels reports which choices of the model option have
// loaded on this device: their files were downloaded then.
func (e *Engine) DownloadedModels() map[string]bool {
	e.modelsMu.Lock()
	defer e.modelsMu.Unlock()
	e.readModels()
	out := map[string]bool{}
	for k, v := range e.models {
		out[k] = v
	}
	return out
}

func (e *Engine) readModels() {
	if e.models != nil {
		return
	}
	e.models = map[string]bool{}
	if data, err := os.ReadFile(e.modelsPath()); err == nil {
		var ids []string
		json.Unmarshal(data, &ids)
		for _, id := range ids {
			e.models[id] = true
		}
	}
}

func (e *Engine) recordModel(id string) {
	if id == "" {
		return
	}
	e.modelsMu.Lock()
	defer e.modelsMu.Unlock()
	e.readModels()
	if e.models[id] {
		return
	}
	e.models[id] = true
	ids := make([]string, 0, len(e.models))
	for k := range e.models {
		ids = append(ids, k)
	}
	sort.Strings(ids)
	data, _ := json.Marshal(ids)
	os.WriteFile(e.modelsPath(), data, 0o644)
}

// InstallModel installs the engine if needed, then loads the model once
// so its files download, and unloads it to give the memory back.
func (e *Engine) InstallModel(ctx context.Context, model string, log engines.Logger) error {
	o := e.m.ModelOption()
	if o == nil {
		return errors.New(e.m.Name + " has no model to choose")
	}
	if !e.Installed() {
		if err := e.Install(ctx, log); err != nil {
			return err
		}
	}
	opts := engines.Resolve(e.m.Options, engines.Values{o.ID: model})
	log("Downloading the model…")
	if _, err := e.ensure(ctx, engines.Subset(e.m.Options, opts, engines.ScopeLoad), log); err != nil {
		return err
	}
	e.Close()
	log("The model is downloaded.")
	return nil
}
