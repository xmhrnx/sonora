package studio

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/tr"
)

// engine returns the selected engine, falling back to the first ready
// one, or nil when there are no voice engines at all.
func (p *Page) engine() *manager.Runtime {
	if rt := p.c.Manager.Get(p.c.Settings.Get().Engine); rt != nil && !rt.Manifest().Transcribes() {
		return rt
	}
	if ready := p.c.Manager.ReadyFor(engines.TaskSpeak); len(ready) > 0 {
		return ready[0]
	}
	if all := p.c.Manager.EnginesFor(engines.TaskSpeak); len(all) > 0 {
		return all[0]
	}
	return nil
}

func (p *Page) selectEngine(id string) {
	p.c.Settings.Update(func(s *settings.Settings) { s.Engine = id })
	if rt := p.c.Manager.Get(id); rt != nil {
		p.c.Manager.LoadVoices(rt, false)
	}
}

// useProfile clones the voice profile with id: it switches to the Clone
// tab, or to an engine that can clone when the current one can't, and
// shows the Studio.
func (p *Page) useProfile(id string) {
	p.c.Settings.Update(func(s *settings.Settings) { s.Profile = id })
	rt := p.engine()
	if rt == nil || rt.Manifest().Clone == engines.CloneNone {
		rt = p.cloningEngine()
		if rt == nil {
			p.c.Info(tr.T("Install an engine that clones voices, such as ZipVoice, to use this profile."))
			p.c.ShowEngines(engines.TaskSpeak)
			return
		}
		p.selectEngine(rt.Manifest().ID)
	}
	p.form(rt).cloneMode = true
	p.c.Go(core.PageStudio)
}

// cloningEngine returns an engine that can clone, ready ones first.
func (p *Page) cloningEngine() *manager.Runtime {
	for _, list := range [][]*manager.Runtime{p.c.Manager.ReadyFor(engines.TaskSpeak), p.c.Manager.EnginesFor(engines.TaskSpeak)} {
		for _, rt := range list {
			if rt.Manifest().Clone != engines.CloneNone {
				return rt
			}
		}
	}
	return nil
}
