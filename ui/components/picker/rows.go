package picker

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

func (p *Picker) list() giode.Element {
	q := strings.ToLower(strings.TrimSpace(p.search.Text()))
	var rows []giode.Element
	group := "\x00"
	for _, it := range p.items() {
		if !it.matches(q) {
			continue
		}
		if it.Group != group {
			group = it.Group
			if group != "" {
				rows = append(rows, giode.Box(giode.Styles{Padding: giode.Inset{Top: 10, Bottom: 4, Left: 10}},
					giode.Text(group, theme.Overline)))
			}
		}
		rows = append(rows, p.row(it))
	}
	if len(rows) == 0 {
		rows = append(rows, giode.Box(giode.Styles{Width: -1, Padding: giode.UniformInset(24)},
			giode.Text(tr.T("Nothing matches “{query}”", "query", q), giode.Styles{Color: "theme.muted", TextAlign: giode.TextAlignCenter, Width: -1})))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 2, Padding: giode.Inset{Right: 10}}, rows...)
}

func (p *Picker) row(it Item) giode.Element {
	return p.rows.Get(it.ID+"\x00"+it.Label+"\x00"+it.Detail, func() *giode.PressableElement {
		return giode.Pressable(func(s giode.PressState) giode.Element {
			sel := p.selected() == it.ID
			bg := ""
			switch {
			case sel:
				bg = "theme.primary/0.16"
			case s.Hovered:
				bg = "theme.text/0.05"
			}
			text := []giode.Element{giode.Text(it.Label, giode.Styles{Color: "theme.text", FontSize: 14, FontWeight: giode.FontWeightMedium})}
			if it.Detail != "" {
				text = append(text, giode.Text(it.Detail, theme.Caption))
			}
			trail := giode.Box(giode.Styles{Width: 18})
			if sel {
				trail = giode.Box(giode.Styles{}, icons.Icon(icons.Check, 18, "theme.primary"))
			}
			return giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 10, Wrap: giode.WrapNoWrap,
				Padding: giode.SymmetricInset(12, 9), BorderRadius: theme.RadiusSm + 2, Background: bg},
				giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2}, text...),
				trail,
			)
		}).Label(it.Label).OnClick(func(giode.ClickEvent) {
			p.open = false
			p.onPick(it)
		})
	})
}
