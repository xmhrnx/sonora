package studio

import (
	"fmt"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/badge"
	"github.com/mehran9675/sonora/ui/components/group"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/statuspill"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// voicePanel is the right-hand column: engine, voice, settings and the
// Generate button.
type voicePanel struct {
	p        *Page
	scroll   *scroll.Scroll
	engine   *picker.Trigger
	install  *giode.ButtonElement
	manage   *giode.ButtonElement
	generate *giode.ButtonElement
	eject    *giode.ButtonElement
	rt       *manager.Runtime // engine shown this frame
}

func newVoicePanel(p *Page) *voicePanel {
	v := &voicePanel{p: p}
	c := p.c
	v.engine = picker.NewTrigger(tr.Msg("Engine"), func() (string, string) {
		label, _ := statuspill.Describe(v.rt)
		return v.rt.Manifest().Name, format.Join(format.Parts(v.rt.Manifest().Vendor), tr.T(label))
	}, func() giode.Element { return badge.Engine(v.rt.Manifest(), 32) }, p.engines.Open)
	v.install = giode.Button("", giode.Styles{BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}).
		Icon(icons.Icon(icons.Download, 15, "")).OnClick(func() {
		rt := v.rt
		done := func(err error) {
			if err == nil {
				c.Success(tr.T("{name} is ready", "name", rt.Manifest().Name))
				c.Manager.LoadVoices(rt, false)
			} else {
				c.Error(err)
			}
		}
		c.Setup(rt, done)
	})
	// callout sets Install's text every frame (core.SetupLabel).
	v.manage = giode.Button("", giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border",
		FontSize: 13, Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}).
		Variant(giode.VariantOutlined).OnClick(func() { c.Go(core.PageEngines) })
	tr.Bind(func(s string) { v.manage.Label(s) }, "Open Engines")
	// footer sets Generate's text every frame.
	v.generate = giode.Button("", giode.Styles{}).OnClick(p.toggleGenerate)
	v.eject = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", FontSize: 12,
		Padding: giode.SymmetricInset(8, 4), BorderRadius: theme.RadiusSm}).Variant(giode.VariantGhost).
		Icon(icons.Icon(icons.Eject, 13, "")).OnClick(func() {
		rt := v.rt
		if err := c.Manager.Eject(rt); err != nil {
			c.Error(err)
			return
		}
		c.Info(tr.T("{name} unloaded — memory freed", "name", rt.Manifest().Name))
	})
	tr.Bind(func(s string) { v.eject.Label(s) }, "Eject")
	v.scroll = giode.Scroll(v.content, theme.Scrolling(5))
	return v
}

func (v *voicePanel) element(rt *manager.Runtime) giode.Element {
	v.rt = rt
	card := theme.Card()
	card.Width, card.Height, card.Gap = 340, -1, 0
	card.Padding = giode.UniformInset(0)
	return giode.Stack(card,
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, v.scroll),
		giode.Divider(giode.Styles{Color: "theme.border"}),
		v.footer(),
	)
}

// compactElement is the same panel for the phone layout: full width,
// no inner scrolling (the whole page scrolls).
func (v *voicePanel) compactElement(rt *manager.Runtime) giode.Element {
	v.rt = rt
	card := theme.Card()
	card.Width, card.Gap = -1, 0
	card.Padding = giode.UniformInset(0)
	return giode.Stack(card, v.content(), giode.Divider(giode.Styles{Color: "theme.border"}), v.footer())
}

func (v *voicePanel) content() giode.Element {
	rt := v.rt
	m := rt.Manifest()
	f := v.p.form(rt)
	els := []giode.Element{group.New(tr.T("Engine"), v.engineSection())}
	if callout := v.callout(); callout != nil {
		els = append(els, callout)
	}
	// The script card shows the hint on wide screens; a phone shows it
	// here.
	if giode.Compact() && m.Hint != "" && rt.State != manager.StateUnavailable {
		els = append(els, giode.Text(tr.T(m.Hint), giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 17}))
	}
	ready := rt.State == manager.StateReady || rt.State == manager.StateWorking
	if ready {
		els = append(els, group.New(tr.T("Voice"), voiceChoice(v.p.cloner, f)))
		if m.Speed {
			els = append(els, giode.Stack(giode.Styles{Width: -1, Gap: 8},
				giode.Row(giode.Styles{Width: -1, Wrap: giode.WrapNoWrap},
					giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("Speed"), theme.Overline)),
					giode.Text(fmt.Sprintf("%.2f×", f.speedValue()), giode.Styles{Color: "theme.text", FontSize: 12, FontWeight: giode.FontWeightBold})),
				syncSlider(f),
			))
		}
		for _, fl := range f.fields {
			els = append(els, fl.Element())
		}
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 22, Padding: giode.UniformInset(20)}, els...)
}
