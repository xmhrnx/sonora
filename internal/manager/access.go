package manager

import (
	"github.com/mehran9675/sonora/internal/engines"
)

// Engines returns every engine in display order.
func (m *Manager) Engines() []*Runtime { return m.list }

// Get returns the engine with id, or nil.
func (m *Manager) Get(id string) *Runtime { return m.byID[id] }

// Shutdown stops every engine's helper process.
func (m *Manager) Shutdown() {
	for _, rt := range m.list {
		if rt.cancel != nil {
			rt.cancel()
		}
		rt.Engine.Close()
	}
}

// CloseIdle stops the helper processes of engines that aren't busy, so
// they restart with new settings on next use.
func (m *Manager) CloseIdle() {
	for _, rt := range m.list {
		if !rt.Busy() {
			rt.Engine.Close()
		}
	}
}

func (m *Manager) environment() engines.Environment {
	s := m.settings.Get()
	return engines.Environment{Device: s.Device, HFToken: s.HFToken}
}

func (m *Manager) manifest(id string) *engines.Manifest {
	if rt := m.byID[id]; rt != nil {
		return rt.Manifest()
	}
	return &engines.Manifest{}
}

// withOptions refreshes a stored custom manifest's option definitions
// from the current template, keeping the user's name.
func withOptions(stored *engines.Manifest, tmpl func(id, name string) *engines.Manifest) *engines.Manifest {
	return tmpl(stored.ID, stored.Name)
}

func initialState(e engines.Engine) State {
	switch {
	case e.Installed():
		return StateReady
	case e.Manifest().NeedsInstall():
		return StateNotInstalled
	default:
		return StateUnavailable
	}
}
