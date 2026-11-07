package transcribe

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/settings"
)

// engine returns the selected speech-to-text engine, falling back to
// the first ready one, or nil when there are none.
func (p *Page) engine() *manager.Runtime {
	m := p.c.Manager
	if rt := m.Get(p.c.Settings.Get().TranscribeEngine); rt != nil && rt.Manifest().Transcribes() {
		return rt
	}
	if ready := m.ReadyFor(engines.TaskTranscribe); len(ready) > 0 {
		return ready[0]
	}
	if all := m.EnginesFor(engines.TaskTranscribe); len(all) > 0 {
		return all[0]
	}
	return nil
}

func (p *Page) selectEngine(id string) {
	p.c.Settings.Update(func(s *settings.Settings) { s.TranscribeEngine = id })
}

// language returns want if the engine understands it, else "" (detect).
func language(m *engines.Manifest, want string) string {
	if len(m.Languages) == 0 {
		return want
	}
	for _, l := range m.Languages {
		if l == want {
			return want
		}
	}
	return ""
}
