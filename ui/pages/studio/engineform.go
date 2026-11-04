package studio

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/slider"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/optionfield"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/segmented"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/tr"
)

// engineForm holds the per-engine controls of the voice panel.
type engineForm struct {
	c  *core.Ctx
	rt *manager.Runtime

	fields       []*optionfield.Field
	voices       *picker.Picker
	voiceTrigger *picker.Trigger
	speed        *slider.Slider

	// cloneMode is set when the user chose "Clone a voice".
	cloneMode bool
	modes     *segmented.Control
}

func newEngineForm(c *core.Ctx, rt *manager.Runtime) *engineForm {
	f := &engineForm{c: c, rt: rt}
	m := rt.Manifest()
	store := c.OptionStore(m, func(id string) {
		if o := m.Option(id); o != nil && o.Scope == engines.ScopeLoad && len(m.Voices) == 0 {
			rt.VoicesLoaded = false
		}
	})
	for _, o := range m.Options {
		if o.Scope != engines.ScopeSetup {
			f.fields = append(f.fields, optionfield.New(o, store))
		}
	}

	f.voices = picker.New(tr.Msg("Choose a voice"), f.voiceItems, func() string { return f.voice().ID }, func(it picker.Item) {
		c.Settings.Update(func(s *settings.Settings) { s.Voices[m.ID] = it.ID })
	})
	f.voiceTrigger = picker.NewTrigger(tr.Msg("Voice"), func() (string, string) {
		if rt.VoicesBusy {
			return tr.T("Loading voices…"), tr.T("Starting {name}", "name", m.Name)
		}
		v := f.voice()
		detail := format.Parts(v.Group)
		if v.Description != "" {
			detail = format.Join(format.Parts(v.Group), format.Parts(v.Description))
		}
		return v.Label(), detail
	}, func() giode.Element {
		return giode.Box(giode.Styles{Width: 30, Height: 30, BorderRadius: 99, Background: "theme.primary/0.15",
			Align: giode.AlignCenter, Justify: giode.JustifyCenter}, icons.Icon(icons.Mic, 15, "theme.primary"))
	}, func() {
		if !rt.VoicesBusy {
			f.voices.Open()
		}
	})

	f.speed = giode.Input.Slider(giode.Styles{Color: "theme.primary", Background: "theme.text/0.12", BorderColor: "theme.primary"}).
		OnChange(func(v float32) {
			c.Settings.Update(func(s *settings.Settings) { s.Speed[m.ID] = sliderToSpeed(v) })
		})
	f.speed.SetValue(speedToSlider(f.speedValue()))

	f.cloneMode = m.Clone == engines.CloneRequired
	f.modes = segmented.New([]segmented.Option{{Value: "voice", Label: tr.Msg("Built-in voice")}, {Value: "clone", Label: tr.Msg("Clone a voice")}},
		func() string {
			if f.cloneMode {
				return "clone"
			}
			return "voice"
		},
		func(v string) { f.cloneMode = v == "clone" })
	return f
}

// modelField returns the field of the engine's model option, or nil.
func (f *engineForm) modelField() *optionfield.Field {
	o := f.rt.Manifest().ModelOption()
	if o == nil {
		return nil
	}
	for _, fl := range f.fields {
		if fl.ID() == o.ID {
			return fl
		}
	}
	return nil
}

func (f *engineForm) voiceList() []engines.Voice {
	if len(f.rt.Voices) > 0 {
		return f.rt.Voices
	}
	return f.rt.Manifest().Voices
}

func (f *engineForm) voiceItems() []picker.Item {
	var out []picker.Item
	for _, v := range f.voiceList() {
		out = append(out, picker.Item{ID: v.ID, Label: v.Label(), Group: format.Parts(v.Group), Detail: format.Parts(v.Description)})
	}
	return out
}

// voice returns the chosen voice, defaulting sensibly.
func (f *engineForm) voice() engines.Voice {
	vs := f.voiceList()
	want := f.c.Settings.Get().Voices[f.rt.Manifest().ID]
	if want == "" {
		want = f.rt.Manifest().DefaultVoice
	}
	for _, v := range vs {
		if v.ID == want {
			return v
		}
	}
	if len(vs) > 0 {
		return vs[0]
	}
	return engines.Voice{ID: want, Name: orDefault(want)}
}

func (f *engineForm) speedValue() float64 {
	if s, ok := f.c.Settings.Get().Speed[f.rt.Manifest().ID]; ok && s > 0 {
		return s
	}
	return 1
}

func sliderToSpeed(v float32) float64 { return float64(int((0.5+float64(v)*1.5)*20+0.5)) / 20 }
func speedToSlider(s float64) float32 { return float32((s - 0.5) / 1.5) }

func orDefault(s string) string {
	if s == "" {
		return tr.T("Default voice")
	}
	return s
}
