package manager

import (
	"fmt"
	"time"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/settings"
)

// AddCustom stores a new API or command engine and returns its id.
func (m *Manager) AddCustom(kind engines.Kind, name string, values engines.Values) string {
	id := fmt.Sprintf("%s-%d", kind, time.Now().UnixNano())
	m.settings.Update(func(s *settings.Settings) {
		s.Custom = append(s.Custom, &engines.Manifest{ID: id, Name: name, Kind: kind, Custom: true})
		s.Values[id] = values
	})
	m.Reload()
	return id
}

// RemoveCustom deletes a user-added engine.
func (m *Manager) RemoveCustom(id string) {
	m.settings.Update(func(s *settings.Settings) {
		for i, c := range s.Custom {
			if c.ID == id {
				s.Custom = append(s.Custom[:i], s.Custom[i+1:]...)
				break
			}
		}
		delete(s.Values, id)
	})
	m.Reload()
}
