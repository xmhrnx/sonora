package studio

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// engineSection is the engine picker, plus an Eject control while the
// engine's model is held in memory.
func (v *voicePanel) engineSection() giode.Element {
	if !v.rt.Engine.Loaded() || v.rt.Busy() {
		return v.engine.Element()
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 6},
		v.engine.Element(),
		giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.Inset{Left: 4}},
			giode.Box(giode.Styles{Width: 8, Height: 8, BorderRadius: 99, Background: "theme.success"}),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("Model loaded in memory"), theme.Caption)),
			v.eject,
		),
	)
}
