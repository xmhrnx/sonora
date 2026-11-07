package transcribe

import (
	"path/filepath"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/badge"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/statuspill"
	"github.com/mehran9675/sonora/ui/components/voicerecorder"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

func newPanel(p *Page) *panel {
	v := &panel{p: p, forms: map[string]*form{}}
	c := p.c
	v.engines = picker.New(tr.Msg("Choose a transcription engine"), func() []picker.Item {
		var ready, other []picker.Item
		for _, rt := range c.Manager.EnginesFor(engines.TaskTranscribe) {
			m := rt.Manifest()
			label, _ := statuspill.Describe(rt)
			it := picker.Item{ID: m.ID, Label: m.Name, Detail: format.Join(format.Parts(m.Vendor), label)}
			if rt.State == manager.StateReady || rt.State == manager.StateWorking {
				it.Group = tr.T("Ready to use")
				ready = append(ready, it)
			} else {
				it.Group = tr.T("Available to install")
				other = append(other, it)
			}
		}
		return append(ready, other...)
	}, func() string { return c.Settings.Get().TranscribeEngine }, func(it picker.Item) { p.selectEngine(it.ID) })
	v.engine = picker.NewTrigger(tr.Msg("Engine"), func() (string, string) {
		label, _ := statuspill.Describe(v.rt)
		return v.rt.Manifest().Name, format.Join(format.Parts(v.rt.Manifest().Vendor), label)
	}, func() giode.Element { return badge.Engine(v.rt.Manifest(), 32) }, v.engines.Open)

	primary := giode.Styles{BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}
	outline := giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border",
		FontSize: 13, Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}
	v.install = giode.Button("", primary).Icon(icons.Icon(icons.Download, 15, "")).OnClick(func() {
		rt := v.rt
		done := func(err error) {
			if err != nil {
				c.Error(err)
				return
			}
			c.Success(tr.T("{name} is ready", "name", rt.Manifest().Name))
		}
		c.Setup(rt, done)
	})
	// callout sets Install's text every frame (core.SetupLabel).
	v.manage = giode.Button("", outline).Variant(giode.VariantOutlined).
		OnClick(func() { c.ShowEngines(engines.TaskTranscribe) })
	tr.Bind(func(s string) { v.manage.Label(s) }, "Open Engines")
	v.pick = giode.Button("", outline).Variant(giode.VariantOutlined).
		Icon(icons.Icon(icons.Upload, 15, "")).OnClick(v.pickFile)
	tr.Bind(func(s string) { v.pick.Label(s) }, "Choose file…")
	v.clear = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(5), BorderRadius: 8}).
		Variant(giode.VariantGhost).Icon(icons.Icon(icons.Close, 14, "")).OnClick(p.dropSource)
	v.recorder = voicerecorder.NewWith(c, voicerecorder.Options{
		Min: time.Second, Name: "speech",
		Hint: tr.Msg("Speak normally; press Stop when you're done."), Ready: tr.Msg("Press Stop when you're done."),
	}, func(path string) {
		p.setSource(path, tr.T("Recording {time}", "time", time.Now().Format("Jan 2, 15:04")), false)
	})
	v.play = giode.Pressable(func(s giode.PressState) giode.Element {
		icon := icons.Solid(icons.Play, 12, "theme.primary")
		if c.IsPlaying(p.source) {
			icon = icons.Solid(icons.Pause, 12, "theme.primary")
		}
		bg := "theme.primary/0.14"
		if s.Hovered {
			bg = "theme.primary/0.26"
		}
		return giode.Box(giode.Styles{Width: 28, Height: 28, BorderRadius: 99, Background: bg,
			Align: giode.AlignCenter, Justify: giode.JustifyCenter}, icon)
	}).OnClick(func(giode.ClickEvent) {
		if filepath.Ext(p.source) != ".wav" {
			c.Info(tr.T("Only WAV files play here; the transcript will play once it's done"))
			return
		}
		c.PlayPause(p.source, p.sourceName, tr.T("Audio to transcribe"))
	})
	tr.Bind(func(s string) { v.play.Label(s) }, "Play audio")
	v.run = giode.Button("", giode.Styles{}).OnClick(p.start)
	v.eject = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", FontSize: 12,
		Padding: giode.SymmetricInset(8, 4), BorderRadius: theme.RadiusSm}).Variant(giode.VariantGhost).
		Icon(icons.Icon(icons.Eject, 13, "")).OnClick(func() {
		if err := c.Manager.Eject(v.rt); err != nil {
			c.Error(err)
			return
		}
		c.Info(tr.T("{name} unloaded — memory freed", "name", v.rt.Manifest().Name))
	})
	tr.Bind(func(s string) { v.eject.Label(s) }, "Eject")
	v.scroll = giode.Scroll(v.content, theme.Scrolling(5))
	return v
}
