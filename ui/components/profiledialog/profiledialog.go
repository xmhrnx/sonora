// Package profiledialog records a new voice profile: a name, a passage
// to read aloud and the recorder. The passage becomes the profile's
// transcript, so cloning engines get the words without anyone typing
// them.
package profiledialog

import (
	"os"
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/components/overlay"
	"github.com/mehran9675/sonora/ui/components/voicerecorder"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Dialog is stateful; create one and Open it as often as needed.
type Dialog struct {
	c       *core.Ctx
	onSaved func(profiles.Profile)

	open    bool
	passage int
	take    string // the recording, until saved or discarded
	saving  bool

	name     *input.Input
	another  *giode.ButtonElement
	recorder *voicerecorder.Recorder
	play     *giode.PressableElement
	redo     *giode.ButtonElement
	cancel   *giode.ButtonElement
	save     *giode.ButtonElement
	close    *giode.ButtonElement
}

// New returns a closed dialog that calls onSaved with each profile it
// saves, on the UI goroutine.
func New(c *core.Ctx, onSaved func(profiles.Profile)) *Dialog {
	d := &Dialog{c: c, onSaved: onSaved}
	d.name = giode.Input.Text("", theme.Field())
	tr.Bind(func(s string) { d.name.Placeholder(s) }, "Name, like “My voice”")
	tr.Bind(func(s string) { d.name.Label(s) }, "Profile name")
	ghost := giode.Styles{Color: "theme.primary", Background: "theme.primary", FontSize: 12,
		Padding: giode.SymmetricInset(8, 5), BorderRadius: theme.RadiusSm}
	d.another = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Refresh, 13, "")).
		OnClick(func() { d.passage = (d.passage + 1) % len(profiles.Passages) })
	tr.Bind(func(s string) { d.another.Label(s) }, "Another text")
	d.recorder = voicerecorder.NewWith(c, voicerecorder.Options{
		Max: profiles.MaxLength, Name: "profile",
		Hint:  tr.Msg("Read the text above at your normal pace."),
		Ready: tr.Msg("Sounds good — press Stop when you've read it all."),
	}, func(path string) { d.discard(); d.take = path })
	d.play = giode.Pressable(func(s giode.PressState) giode.Element {
		icon := icons.Play
		if c.IsPlaying(d.take) {
			icon = icons.Pause
		}
		return giode.Box(giode.Styles{Width: 34, Height: 34, BorderRadius: 99, Background: "theme.primary/0.18",
			Align: giode.AlignCenter, Justify: giode.JustifyCenter}, icons.Solid(icon, 14, "theme.primary"))
	}).OnClick(func(giode.ClickEvent) { c.PlayPause(d.take, tr.T("Voice sample"), d.profileName()) })
	tr.Bind(func(s string) { d.play.Label(s) }, "Play recording")
	outline := giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border", FontSize: 13,
		Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}
	d.redo = giode.Button("", outline).Variant(giode.VariantOutlined).Icon(icons.Icon(icons.Mic, 14, "")).
		OnClick(func() { d.discard() })
	tr.Bind(func(s string) { d.redo.Label(s) }, "Record again")
	d.cancel = giode.Button("", outline).Variant(giode.VariantOutlined).OnClick(d.Close)
	tr.Bind(func(s string) { d.cancel.Label(s) }, "Cancel")
	d.save = giode.Button("", giode.Styles{BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(16, 8), BorderRadius: theme.RadiusSm}).
		Icon(icons.Icon(icons.Check, 15, "")).OnClick(d.saveProfile)
	tr.Bind(func(s string) { d.save.Label(s) }, "Save profile")
	d.close = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(6), BorderRadius: 8}).
		Variant(giode.VariantGhost).Icon(icons.Icon(icons.Close, 16, "")).OnClick(d.Close)
	overlay.Add(giode.Dialog(&d.open, d.content, theme.Dialog()))
	return d
}

// Open shows the dialog with an empty name and the next passage.
func (d *Dialog) Open() {
	d.discard()
	d.name.SetText("")
	d.passage = (d.passage + 1) % len(profiles.Passages)
	d.open = true
}

// Close discards an unsaved take and closes the dialog.
func (d *Dialog) Close() {
	if d.recorder.Recording() {
		return // Stop or Cancel the recording first
	}
	d.discard()
	d.open = false
}

// discard deletes the take, if any.
func (d *Dialog) discard() {
	if d.take != "" {
		os.Remove(d.take)
		d.take = ""
	}
}

// profileName is the name typed, or a default.
func (d *Dialog) profileName() string {
	if n := strings.TrimSpace(d.name.Text()); n != "" {
		return n
	}
	return tr.T("My voice")
}

// saveProfile stores the take with the passage as its transcript.
func (d *Dialog) saveProfile() {
	if d.take == "" || d.saving {
		return
	}
	d.saving = true
	name, passage, take := d.profileName(), profiles.Passages[d.passage], d.take
	var p profiles.Profile
	d.c.Background(func() error {
		var err error
		p, err = d.c.Profiles.Add(name, passage, take, false)
		return err
	}, func(err error) {
		d.saving = false
		if err != nil {
			d.c.Error(err)
			return
		}
		d.discard()
		d.open = false
		d.c.Success(tr.T("Saved the profile {name}", "name", p.Name))
		if d.onSaved != nil {
			d.onSaved(p)
		}
	})
}

func (d *Dialog) content() giode.Element {
	title := giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("New voice profile"), theme.Heading)), d.close)
	passage := giode.Stack(giode.Styles{Width: -1, Gap: 8, Padding: giode.UniformInset(14),
		Background: "theme.primary/0.08", BorderRadius: theme.RadiusMd, BorderWidth: 1, BorderColor: "theme.primary/0.25"},
		giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("Read this aloud"), theme.Overline)), d.another),
		giode.Text(profiles.Passages[d.passage], giode.Styles{Color: "theme.text", FontSize: 16, LineHeight: 25}),
	)
	d.another.Disabled(d.recorder.Recording() || d.take != "")
	var take giode.Element
	switch {
	case d.recorder.Recording():
		take = d.recorder.Element()
	case d.take != "":
		take = giode.Row(giode.Styles{Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
			Padding: giode.SymmetricInset(10, 8), Background: "theme.primary/0.1", BorderRadius: theme.RadiusSm + 2},
			d.play,
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("Recording ready"), giode.Styles{Color: "theme.text", FontSize: 13})),
			d.redo)
	default:
		take = giode.Row(giode.Styles{Width: -1, Gap: 10, Align: giode.AlignCenter},
			d.recorder.Button(),
			giode.Text(tr.T("Then press Stop. 5 to 30 seconds is enough."), theme.Caption))
	}
	d.save.Disabled(d.take == "").Loading(d.saving)
	d.cancel.Disabled(d.recorder.Recording())
	return giode.Stack(giode.Styles{Width: -1, MaxWidth: 520, Gap: 16},
		title,
		giode.Text(tr.T("Record yourself reading the text below in a quiet room. Sonora gives the engine the recording and the words, so there's nothing to type."),
			giode.Styles{Color: "theme.secondary", FontSize: 13, LineHeight: 19}),
		giode.Box(giode.Styles{Width: -1}, d.name),
		passage,
		take,
		giode.Row(giode.Styles{Width: -1, Gap: 8, Justify: giode.JustifyEnd}, d.cancel, d.save),
	)
}
