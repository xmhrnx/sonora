package manager

import (
	"github.com/mehran9675/sonora/internal/engines"
)

// Ready returns the engines that can work now.
func (m *Manager) Ready() []*Runtime { return ready(m.list) }

// ReadyFor returns the engines that can do task now.
func (m *Manager) ReadyFor(task engines.Task) []*Runtime { return ready(m.EnginesFor(task)) }

func ready(list []*Runtime) []*Runtime {
	var out []*Runtime
	for _, rt := range list {
		if rt.State == StateReady || rt.State == StateWorking {
			out = append(out, rt)
		}
	}
	return out
}

// EnginesFor returns the engines that do task, in display order.
func (m *Manager) EnginesFor(task engines.Task) []*Runtime {
	var out []*Runtime
	for _, rt := range m.list {
		if rt.Manifest().Task == task {
			out = append(out, rt)
		}
	}
	return out
}
