package library

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/playbutton"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// recRow is one recording in the list.
type recRow struct {
	c          *core.Ctx
	r          recording
	play       *giode.PressableElement
	transcribe *giode.ButtonElement
	export     *giode.ButtonElement
	remove     *giode.ButtonElement
}

func newRecRow(c *core.Ctx, l *recordingList, r recording) *recRow {
	row := &recRow{c: c, r: r}
	row.play = playbutton.New(c, r.Path, r.Name, tr.T("Recording"))
	ghost := giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(8), BorderRadius: 8}
	row.transcribe = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Captions, 16, "")).
		OnClick(func() { c.Transcribe(r.Path, r.Name) })
	row.export = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Share, 16, "")).
		OnClick(func() { c.ExportFile(r.Path, r.Name) })
	row.remove = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Trash, 16, "")).OnClick(func() {
		c.Toasts.ShowAction(giode.ToastWarning, tr.T("Delete {name}?", "name", r.Name), tr.T("Delete"), func() {
			giode.AfterLayout(func() {
				if c.Player.Loaded() == r.Path {
					c.Player.Unload()
					c.Now = nil
				}
				if err := os.Remove(r.Path); err != nil {
					c.Error(err)
				}
				l.refresh()
			})
		})
	})
	return row
}

func (row *recRow) element() giode.Element {
	r := row.r
	card := theme.Card()
	card.FlexDirection, card.Align, card.Gap, card.Width = giode.FlexDirectionRow, giode.AlignCenter, 16, -1
	card.Padding = giode.SymmetricInset(16, 14)
	card.Wrap = giode.WrapNoWrap
	meta := []giode.Element{chip.Chip(strings.ToUpper(strings.TrimPrefix(filepath.Ext(r.Path), ".")))}
	info := sizeLabel(r.Size) + "  ·  " + format.Ago(r.Modified)
	if r.Duration > 0 {
		info = format.Clock(r.Duration) + "  ·  " + info
	}
	meta = append(meta, giode.Text(info, theme.Caption))
	return giode.Row(card,
		row.play,
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 8},
			giode.Text(r.Name, giode.Styles{Color: "theme.text", FontSize: 14, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter}, meta...),
		),
		giode.Row(giode.Styles{Gap: 2, Wrap: giode.WrapNoWrap}, row.transcribe, row.export, row.remove),
	)
}

// sizeLabel formats a file size.
func sizeLabel(n int64) string {
	switch {
	case n >= 1<<30:
		return tr.T("{size} GB", "size", fmt.Sprintf("%.1f", float64(n)/(1<<30)))
	case n >= 1<<20:
		return tr.T("{size} MB", "size", fmt.Sprintf("%.1f", float64(n)/(1<<20)))
	}
	return tr.T("{size} KB", "size", (n+1023)>>10)
}
