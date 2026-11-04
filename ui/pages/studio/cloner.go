package studio

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/profiledialog"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// cloner is the Clone tab: the voice profile to copy, chosen from the
// Profiles page's list or recorded right here. The engine gets the
// profile's recording and its transcript; nobody types anything.
type cloner struct {
	c       *core.Ctx
	list    *picker.Picker
	trigger *picker.Trigger
	play    *giode.PressableElement
	create  *giode.ButtonElement
	manage  *giode.ButtonElement
	fill    *giode.ButtonElement
	dialog  *profiledialog.Dialog
	filling bool
}

func newCloner(c *core.Ctx) *cloner {
	cl := &cloner{c: c}
	choose := func(id string) { c.Settings.Update(func(s *settings.Settings) { s.Profile = id }) }
	cl.dialog = profiledialog.New(c, func(p profiles.Profile) { choose(p.ID) })
	cl.list = picker.New(tr.Msg("Choose a voice profile"), func() []picker.Item {
		var items []picker.Item
		for _, p := range c.Profiles.List() {
			items = append(items, picker.Item{ID: p.ID, Label: p.Name, Detail: profileDetail(p)})
		}
		return items
	}, func() string { p, _ := cl.selected(); return p.ID }, func(it picker.Item) { choose(it.ID) })
	cl.trigger = picker.NewTrigger(tr.Msg("Voice profile"), func() (string, string) {
		p, ok := cl.selected()
		if !ok {
			return tr.T("Choose a voice profile"), ""
		}
		return p.Name, profileDetail(p)
	}, nil, cl.list.Open)
	cl.play = giode.Pressable(func(s giode.PressState) giode.Element {
		p, _ := cl.selected()
		icon := icons.Solid(icons.Play, 13, "theme.primary")
		if c.IsPlaying(c.Profiles.Path(p)) {
			icon = icons.Solid(icons.Pause, 13, "theme.primary")
		}
		bg := "theme.primary/0.14"
		if s.Hovered {
			bg = "theme.primary/0.26"
		}
		return giode.Box(giode.Styles{Width: 36, Height: 36, BorderRadius: 99, Background: bg,
			Align: giode.AlignCenter, Justify: giode.JustifyCenter}, icon)
	}).OnClick(func(giode.ClickEvent) {
		if p, ok := cl.selected(); ok {
			c.PlayPause(c.Profiles.Path(p), p.Name, tr.T("Voice profile"))
		}
	})
	tr.Bind(func(s string) { cl.play.Label(s) }, "Play recording")
	ghost := giode.Styles{Color: "theme.primary", Background: "theme.primary", FontSize: 12,
		Padding: giode.SymmetricInset(8, 5), BorderRadius: theme.RadiusSm}
	cl.create = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Plus, 13, "")).OnClick(cl.dialog.Open)
	tr.Bind(func(s string) { cl.create.Label(s) }, "New profile")
	cl.manage = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.UserVoice, 13, "")).
		OnClick(func() { c.Go(core.PageProfiles) })
	tr.Bind(func(s string) { cl.manage.Label(s) }, "All profiles")
	cl.fill = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Captions, 13, "")).OnClick(func() {
		if p, ok := cl.selected(); ok {
			cl.filling = true
			c.FillTranscript(p, func(error) { cl.filling = false })
		}
	})
	tr.Bind(func(s string) { cl.fill.Label(s) }, "Find out what is said")
	return cl
}

// selected is the profile to clone: the one chosen, or the newest.
func (cl *cloner) selected() (profiles.Profile, bool) {
	if p, ok := cl.c.Profiles.Get(cl.c.Settings.Get().Profile); ok {
		return p, true
	}
	if all := cl.c.Profiles.List(); len(all) > 0 {
		return all[0], true
	}
	return profiles.Profile{}, false
}

// profileDetail describes a profile in a list: its length and what is
// said.
func profileDetail(p profiles.Profile) string {
	said := tr.T("No transcript yet")
	if p.Transcript != "" {
		said = "“" + p.Transcript + "”"
	}
	return format.Join(format.Clock(p.Duration), said)
}

// element lays out the tab for an engine described by m.
func (cl *cloner) element(m *engines.Manifest) giode.Element {
	note := giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 17}
	links := giode.Row(giode.Styles{Gap: 4, Wrap: giode.WrapNoWrap}, cl.create, cl.manage)
	p, ok := cl.selected()
	if !ok {
		return giode.Stack(giode.Styles{Width: -1, Gap: 10},
			giode.Text(tr.T("Clone a voice from a profile: record yourself reading a short text, or import a recording of one person. "+
				"Sonora gives the engine the words that are said."), note),
			links)
	}
	body := []giode.Element{
		giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			cl.play, giode.Box(giode.Styles{FlexGrow: 1}, cl.trigger.Element())),
		links,
	}
	if p.Transcript == "" && m.Transcript != engines.CloneNone {
		msg := tr.T("{engine} works out what is said itself.", "engine", m.Name)
		style := note
		if m.Transcript == engines.CloneRequired {
			msg = tr.T("{engine} needs to know what is said in this recording.", "engine", m.Name)
			style.Color = "theme.danger"
		}
		row := []giode.Element{giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(msg, style))}
		switch {
		case cl.filling:
			row = append(row, giode.Text(tr.T("Finding out what is said…"), theme.Caption))
		case cl.c.CanFillTranscripts():
			row = append(row, cl.fill)
		case m.Transcript == engines.CloneRequired:
			row[0] = giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(
				tr.T("{engine} needs to know what is said in this recording. Install a transcription engine, such as Whisper, to fill it in.",
					"engine", m.Name), style))
		}
		body = append(body, giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter}, row...))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10}, body...)
}
