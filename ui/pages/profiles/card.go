package profiles

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/input"

	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/playbutton"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// card is one profile: play, its name (renamed in place), how long it
// is, what is said in it, and Use in Studio, rename and delete.
type card struct {
	c  *core.Ctx
	id string

	play       *giode.PressableElement
	use        *giode.ButtonElement
	rename     *giode.ButtonElement
	remove     *giode.ButtonElement
	transcribe *giode.ButtonElement

	editing bool
	name    *input.Input
	done    *giode.ButtonElement
	filling bool // a transcription engine is listening to it
}

func newCard(c *core.Ctx, id string) *card {
	cd := &card{c: c, id: id}
	pr, _ := c.Profiles.Get(id)
	cd.play = playbutton.New(c, c.Profiles.Path(pr), pr.Name, tr.T("Voice profile"))
	ghost := giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(8), BorderRadius: 8}
	cd.use = giode.Button("", giode.Styles{Color: "theme.primary", Background: "theme.primary/0.12", FontSize: 13,
		Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}).
		Icon(icons.Icon(icons.Sparkles, 14, "")).OnClick(func() {
		if c.UseProfile != nil {
			c.UseProfile(id)
		}
	})
	tr.Bind(func(s string) { cd.use.Label(s) }, "Use in Studio")
	cd.rename = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Pencil, 16, "")).OnClick(func() {
		pr, _ := c.Profiles.Get(id)
		cd.name.SetText(pr.Name)
		cd.name.Focus()
		cd.editing = true
	})
	cd.remove = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Trash, 16, "")).OnClick(func() {
		pr, _ := c.Profiles.Get(id)
		c.Toasts.ShowAction(giode.ToastWarning, tr.T("Delete {name}?", "name", pr.Name), tr.T("Delete"), func() {
			giode.AfterLayout(func() {
				if path := c.Profiles.Path(pr); c.Player.Loaded() == path {
					c.Player.Unload()
					c.Now = nil
				}
				c.Error(c.Profiles.Delete(id))
			})
		})
	})
	cd.transcribe = giode.Button("", giode.Styles{Color: "theme.primary", Background: "theme.primary", FontSize: 12,
		Padding: giode.SymmetricInset(8, 5), BorderRadius: theme.RadiusSm}).Variant(giode.VariantGhost).
		Icon(icons.Icon(icons.Captions, 13, "")).OnClick(func() {
		pr, _ := c.Profiles.Get(id)
		cd.filling = true
		c.FillTranscript(pr, func(error) { cd.filling = false })
	})
	tr.Bind(func(s string) { cd.transcribe.Label(s) }, "Find out what is said")
	cd.name = giode.Input.Text("", theme.Field()).Submit(func(string) { cd.saveName() })
	tr.Bind(func(s string) { cd.name.Label(s) }, "Profile name")
	cd.done = giode.Button("", giode.Styles{Color: "theme.primary", Background: "theme.primary/0.12", FontSize: 13,
		Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}).
		Icon(icons.Icon(icons.Check, 14, "")).OnClick(cd.saveName)
	tr.Bind(func(s string) { cd.done.Label(s) }, "Save")
	return cd
}

func (cd *card) element(pr profiles.Profile) giode.Element {
	st := theme.Card()
	st.FlexDirection, st.Align, st.Gap, st.Width = giode.FlexDirectionRow, giode.AlignCenter, 16, -1
	st.Padding, st.Wrap = giode.SymmetricInset(16, 14), giode.WrapNoWrap
	kind := tr.T("Recorded")
	if pr.Imported {
		kind = tr.T("Imported")
	}
	meta := giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter},
		chip.Chip(format.Clock(pr.Duration)),
		giode.Text(format.Join(kind, format.Ago(pr.Created)), theme.Caption))
	name := giode.Element(giode.Text(pr.Name, giode.Styles{Color: "theme.text", FontSize: 15, FontWeight: giode.FontWeightBold,
		TextOverflow: giode.TextOverflowEllipsis}))
	actions := giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, cd.use, cd.rename, cd.remove)
	if cd.editing {
		name = giode.Box(giode.Styles{Width: -1, MaxWidth: 320}, cd.name)
		actions = giode.Row(giode.Styles{Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, cd.done)
	}
	body := []giode.Element{name, meta, cd.transcript(pr)}
	if giode.Compact() {
		// Phones put the actions under the text.
		body = append(body, actions)
		return giode.Row(st, cd.play, giode.Stack(giode.Styles{FlexGrow: 1, Gap: 8}, body...))
	}
	return giode.Row(st, cd.play, giode.Stack(giode.Styles{FlexGrow: 1, Gap: 8}, body...), actions)
}
