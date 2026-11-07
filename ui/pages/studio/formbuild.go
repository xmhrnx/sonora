package studio

import (
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/settings"
)

// form returns the cached settings form for rt.
func (p *Page) form(rt *manager.Runtime) *engineForm {
	id := rt.Manifest().ID
	f := p.forms[id]
	if f == nil || f.rt != rt {
		f = newEngineForm(p.c, rt)
		p.forms[id] = f
	}
	return f
}

// SaveDraft stores the editor text. It's safe from any goroutine, as
// OnBackground on Android requires.
func (p *Page) SaveDraft() {
	text := p.editor.savedText()
	p.c.Settings.Update(func(s *settings.Settings) { s.Draft = text })
}
