package manager

import (
	"log"
	"path/filepath"

	"github.com/mehran9675/sonora/internal/catalog"
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/command"
	"github.com/mehran9675/sonora/internal/engines/onnx"
	"github.com/mehran9675/sonora/internal/engines/openai"
	"github.com/mehran9675/sonora/internal/engines/python"
	"github.com/mehran9675/sonora/internal/engines/system"
	"github.com/mehran9675/sonora/internal/library"
	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/internal/transcripts"
)

// New unpacks the Python runtime and builds the engine list. post must
// run a function on the UI goroutine (giode's App.Post).
func New(post func(func()), st *settings.Store, lib *library.Library, ts *transcripts.Store) *Manager {
	m := &Manager{post: post, settings: st, library: lib, transcripts: ts, byID: map[string]*Runtime{}}
	if !mobile {
		script, err := catalog.Unpack(paths.Sub("runtime"))
		if err != nil {
			log.Printf("unpack runtime: %v", err)
		}
		m.script = script
	}
	m.Reload()
	return m
}

// Reload rebuilds the engine list (after adding or removing engines),
// keeping the runtimes of engines that still exist.
func (m *Manager) Reload() {
	old := m.byID
	m.list, m.byID = nil, map[string]*Runtime{}
	add := func(e engines.Engine) {
		id := e.Manifest().ID
		if _, dup := m.byID[id]; dup {
			return
		}
		rt := old[id]
		if rt == nil {
			rt = &Runtime{Engine: e, Progress: -1}
			rt.State = initialState(e)
		}
		m.list = append(m.list, rt)
		m.byID[id] = rt
	}
	// On-device engines run inside the app (builds with sherpa-onnx).
	if onnx.Available {
		for _, e := range onnx.All(paths.Data()) {
			add(e)
		}
	}
	// Phones can't run local programs or Python.
	if !mobile {
		add(system.New())
		builtin, errs := catalog.LoadDir(filepath.Join(paths.Data(), "runtime", "engines"))
		plugins, perrs := catalog.LoadDir(paths.Sub("plugins"))
		for _, err := range append(errs, perrs...) {
			log.Printf("engine manifest: %v", err)
		}
		for _, man := range append(builtin, plugins...) {
			add(python.New(man, paths.Data(), m.script, m.environment))
		}
	}
	for _, man := range m.settings.Get().Custom {
		id := man.ID
		cfg := func() engines.Values { return engines.Resolve(m.manifest(id).Options, m.settings.Values(id)) }
		switch man.Kind {
		case engines.KindOpenAI:
			add(openai.New(withOptions(man, openai.Template), cfg))
		case engines.KindOpenAITranscribe:
			add(openai.NewTranscriber(withOptions(man, openai.TranscribeTemplate), cfg))
		case engines.KindCommand:
			if mobile {
				continue
			}
			add(command.New(withOptions(man, command.Template), cfg))
		}
	}
}
