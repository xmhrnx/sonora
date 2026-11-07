package picker

import (
	"runtime"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/mehran9675/sonora/ui/components/overlay"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Picker is a modal list. Create it once; it registers itself as an
// overlay.
type Picker struct {
	open     bool
	title    string
	items    func() []Item
	selected func() string
	onPick   func(Item)

	search *input.Input
	// searchable shows the search field, for lists of searchFrom or more.
	searchable bool
	scroll     *scroll.Scroll
	rows       *giode.KeyedCache[string, *giode.PressableElement]
	close      *giode.ButtonElement
	dialog     *dialog.Dialog
}

// searchFrom is the shortest list that gets a search field.
const searchFrom = 9

// touch is true on phones, where focusing the search field would open
// the keyboard over the list.
const touch = runtime.GOOS == "android" || runtime.GOOS == "ios"

// New returns a picker titled title (English, shown translated)
// listing items(); selected returns the current choice's ID; onPick
// runs when the user chooses.
func New(title string, items func() []Item, selected func() string, onPick func(Item)) *Picker {
	p := &Picker{title: title, items: items, selected: selected, onPick: onPick}
	field := theme.Field()
	field.Padding = giode.SymmetricInset(12, 10)
	p.search = giode.Input.Text("", field)
	tr.Bind(func(s string) { p.search.Placeholder(s) }, "Search…")
	p.rows = giode.NewKeyedCache[string, *giode.PressableElement]()
	p.scroll = giode.Scroll(p.list, theme.Scrolling(6))
	p.close = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(6), BorderRadius: 8}).
		Variant(giode.VariantGhost).Icon(icons.Icon(icons.Close, 16, "")).OnClick(func() { p.open = false })
	p.dialog = giode.Dialog(&p.open, p.content, theme.Dialog())
	overlay.Add(p.dialog)
	return p
}

// Open shows the picker with an empty search, focused on desktops.
func (p *Picker) Open() {
	p.search.SetText("")
	p.searchable = len(p.items()) >= searchFrom
	if p.searchable && !touch {
		p.search.Focus()
	}
	p.scroll.ScrollTo(0)
	p.open = true
}

func (p *Picker) content() giode.Element {
	parts := []giode.Element{
		giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T(p.title), theme.Heading)),
			p.close,
		),
	}
	if p.searchable {
		parts = append(parts, giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 8, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, p.search),
		))
	}
	parts = append(parts, giode.Box(giode.Styles{Width: -1, Height: 380}, p.scroll))
	return giode.Stack(giode.Styles{Width: -1, MaxWidth: 460, Gap: 14}, parts...)
}
