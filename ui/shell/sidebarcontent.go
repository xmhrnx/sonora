package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Tabs lays out the same navigation as a bottom bar for narrow screens.
func (s *Sidebar) Tabs() giode.Element {
	tabs := make([]giode.Element, len(s.items))
	for i, it := range s.items {
		tabs[i] = it.Tab()
	}
	return giode.Row(giode.Styles{Width: -1, Gap: 4, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(8, 6),
		Background: "theme.surface"}, tabs...)
}

// Element lays out the sidebar.
func (s *Sidebar) Element() giode.Element {
	els := []giode.Element{logo(), giode.Box(giode.Styles{Height: 18})}
	for _, it := range s.items {
		els = append(els, it.Element())
	}
	els = append(els,
		giode.Box(giode.Styles{FlexGrow: 1}),
		s.engineSummary(),
		s.theme,
	)
	return giode.Stack(giode.Styles{
		Width: 232, Height: -1, Gap: 4, Padding: giode.Inset{Top: 18, Left: 14, Right: 14, Bottom: 14},
		Background: "theme.surface", BorderWidth: 0,
	}, els...)
}

func (s *Sidebar) engineSummary() giode.Element {
	ready := len(s.c.Manager.Ready())
	msg := tr.N("{n} engines ready", ready)
	m := s.c.Manager
	for _, rt := range m.Engines() {
		if rt.Busy() && rt.State == manager.StateInstalling {
			msg = tr.T("Installing {name}", "name", rt.Manifest().Name)
		}
	}
	if j := m.TranscribeJob; j != nil {
		msg = tr.T("Transcribing with {name}", "name", j.Runtime.Manifest().Name)
	}
	if j := m.Job; j != nil {
		msg = tr.T("Generating with {name}", "name", j.Runtime.Manifest().Name)
	}
	return giode.Row(giode.Styles{Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.UniformInset(12), Margin: giode.Inset{Bottom: 6},
		Background: "theme.text/0.04", BorderRadius: theme.RadiusMd},
		giode.Box(giode.Styles{Width: 8, Height: 8, BorderRadius: 99, Background: "theme.success"}),
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(msg, giode.Styles{Color: "theme.secondary", FontSize: 12, MaxLines: 2})),
	)
}
