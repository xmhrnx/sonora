// Package profiles is the Profiles screen: voices to clone, recorded
// reading a passage aloud or imported from a file, each with what is
// said in it.
package profiles

import (
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/mehran9675/sonora/ui/components/emptystate"
	"github.com/mehran9675/sonora/ui/components/pagehead"
	"github.com/mehran9675/sonora/ui/components/playerbar"
	"github.com/mehran9675/sonora/ui/components/profiledialog"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Page is the Profiles screen. Create it once.
type Page struct {
	c         *core.Ctx
	scroll    *scroll.Scroll
	cards     *giode.KeyedCache[string, *card]
	player    *playerbar.Bar
	dialog    *profiledialog.Dialog
	record    *giode.ButtonElement
	importBtn *giode.ButtonElement
	importing bool
}

// New builds the page.
func New(c *core.Ctx) *Page {
	p := &Page{c: c, cards: giode.NewKeyedCache[string, *card]()}
	p.scroll = giode.Scroll(p.list, theme.Scrolling(6))
	p.player = playerbar.New(c)
	p.dialog = profiledialog.New(c, nil)
	p.record = giode.Button("", giode.Styles{BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(16, 9), BorderRadius: theme.RadiusSm + 2}).
		Icon(icons.Icon(icons.Mic, 15, "")).OnClick(p.dialog.Open)
	tr.Bind(func(s string) { p.record.Label(s) }, "New profile")
	p.importBtn = giode.Button("", giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border",
		FontSize: 13, Padding: giode.SymmetricInset(14, 9), BorderRadius: theme.RadiusSm + 2}).
		Variant(giode.VariantOutlined).Icon(icons.Icon(icons.Upload, 15, "")).OnClick(p.importFile)
	tr.Bind(func(s string) { p.importBtn.Label(s) }, "Import a recording")
	return p
}

func (p *Page) list() giode.Element {
	all := p.c.Profiles.List()
	if len(all) == 0 {
		return emptystate.New(icons.UserVoice, tr.T("No voice profiles yet"),
			tr.T("A profile is a voice to clone: record yourself reading a short text, or import a recording of one person speaking. "+
				"Then pick it on the Studio's Clone tab."),
			p.record, p.importBtn)
	}
	els := make([]giode.Element, 0, len(all))
	for _, pr := range all {
		els = append(els, p.cards.Get(pr.ID, func() *card { return newCard(p.c, pr.ID) }).element(pr))
	}
	// Room for the scrollbar, which phones don't show.
	right := 12
	if giode.Compact() {
		right = 0
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10, Padding: giode.Inset{Right: right, Bottom: 12}}, els...)
}

// Element lays out the page.
func (p *Page) Element() giode.Element {
	p.importBtn.Loading(p.importing)
	var actions []giode.Element
	if len(p.c.Profiles.List()) > 0 {
		actions = append(actions, giode.Row(giode.Styles{Gap: 8}, p.importBtn, p.record))
	}
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Gap: 20, Padding: giode.UniformInset(core.Gutter())},
		pagehead.Head(giode.Compact(), tr.T("Profiles"), tr.T("Voices to clone, each with what is said in it."), actions...),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, p.scroll),
		p.player.Element(),
	)
}

// title turns a file name into a profile name: "my-voice_01.wav" reads
// "my voice 01".
func title(name string) string {
	if i := strings.LastIndexByte(name, '.'); i > 0 {
		name = name[:i]
	}
	name = strings.NewReplacer("_", " ", "-", " ").Replace(name)
	if name = strings.TrimSpace(name); name == "" {
		return tr.T("Imported voice")
	}
	return name
}
