package engines

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/ui/components/badge"
	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/logdialog"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/statuspill"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// card is one engine tile.
type card struct {
	c  *core.Ctx
	rt *manager.Runtime

	install   *giode.ButtonElement
	cancel    *giode.ButtonElement
	use       *giode.ButtonElement
	logBtn    *giode.ButtonElement
	remove    *giode.ButtonElement
	configure *giode.ButtonElement
	site      *giode.ButtonElement
	eject     *giode.ButtonElement

	// Choosing a model to download, for engines with a model choice:
	// choosing shows the chooser in place of the card's details.
	choosing bool
	choice   string
	models   *picker.Picker
	model    *picker.Trigger
	download *giode.ButtonElement
	back     *giode.ButtonElement
	addModel *giode.ButtonElement
}

func newCard(c *core.Ctx, rt *manager.Runtime, cfg *configDialog, logs *logdialog.Dialog) *card {
	cd := &card{c: c, rt: rt}
	m := rt.Manifest()
	primary := giode.Styles{BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(14, 8), BorderRadius: theme.RadiusSm}
	outline := giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border", FontSize: 13,
		Padding: giode.SymmetricInset(12, 8), BorderRadius: theme.RadiusSm}
	ghost := giode.Styles{Color: "theme.muted", Background: "theme.text", FontSize: 13,
		Padding: giode.SymmetricInset(10, 8), BorderRadius: theme.RadiusSm}

	cd.install = giode.Button("", primary).Icon(icons.Icon(icons.Download, 15, "")).OnClick(func() {
		// Engines that are their models download the one the user picks.
		if len(engines.Models(rt.Engine)) > 0 {
			cd.choose()
			return
		}
		c.Install(rt, func(err error) {
			if err != nil {
				c.Error(err)
				logs.Show(rt)
				return
			}
			c.Success(tr.T("{engine} is installed", "engine", m.Name))
		})
	})
	tr.Bind(func(s string) { cd.install.Label(s) }, "Install")
	if len(engines.Models(rt.Engine)) > 0 {
		cd.initChooser(primary, outline, logs)
	}
	cd.cancel = giode.Button("", outline).Variant(giode.VariantOutlined).OnClick(func() { c.Manager.Cancel(rt) })
	tr.Bind(func(s string) { cd.cancel.Label(s) }, "Cancel")
	cd.use = giode.Button("", primary).Icon(icons.Icon(icons.Sparkles, 15, "")).OnClick(func() {
		c.UseEngine(m.ID)
	})
	use := tr.Msg("Use in Studio")
	if m.Transcribes() {
		cd.use = giode.Button("", primary).Icon(icons.Icon(icons.Captions, 15, "")).OnClick(func() {
			c.UseTranscriber(m.ID)
		})
		use = tr.Msg("Transcribe")
	}
	tr.Bind(func(s string) { cd.use.Label(s) }, use)
	cd.logBtn = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Terminal, 15, "")).
		OnClick(func() { logs.Show(rt) })
	cd.remove = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Trash, 15, "")).OnClick(func() {
		if m.Custom {
			c.Manager.RemoveCustom(m.ID)
			c.Info(tr.T("Removed {engine}", "engine", m.Name))
			return
		}
		c.Toasts.ShowAction(giode.ToastWarning, tr.T("Remove {engine} and its downloads?", "engine", m.Name), tr.T("Remove"), func() {
			if err := c.Manager.Uninstall(rt); err != nil {
				c.Error(err)
				return
			}
			c.Info(tr.T("{engine} removed", "engine", m.Name))
		})
	})
	cd.configure = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Sliders, 15, "")).
		OnClick(func() { cfg.show(rt) })
	cd.site = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.External, 15, "")).OnClick(func() {
		c.Error(paths.Open(m.Homepage))
	})
	cd.eject = giode.Button("", ghost).Variant(giode.VariantGhost).Icon(icons.Icon(icons.Eject, 15, "")).OnClick(func() {
		if err := c.Manager.Eject(rt); err != nil {
			c.Error(err)
			return
		}
		c.Info(tr.T("{engine} unloaded from memory", "engine", m.Name))
	})
	return cd
}

// Fixed heights for each part of a card, so every card is the same
// size whatever its state.
const (
	cardHeight   = 356
	descHeight   = 64 // three lines
	metaHeight   = 32 // two lines
	statusHeight = 38
)

func (cd *card) element() giode.Element {
	rt := cd.rt
	m := rt.Manifest()
	st := theme.Card()
	st.Width, st.Height, st.Gap, st.Overflow = 318, cardHeight, 12, giode.OverflowHidden
	if giode.Compact() {
		st.Width = -1 // one card per row on a phone
	}
	header := giode.Row(giode.Styles{Width: -1, Height: 48, Gap: 12, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		badge.Engine(m, 44),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 4},
			giode.Text(m.Name, giode.Styles{Color: "theme.text", FontSize: 16, FontWeight: giode.FontWeightBold, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				statuspill.For(rt),
				giode.Text(format.Parts(m.Vendor), giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis}),
			),
		),
	)
	var tags []giode.Element
	for _, t := range m.Tags {
		tags = append(tags, chip.Chip(tr.T(t)))
	}
	if cd.choosing && cd.rt.State != manager.StateInstalling {
		// The chooser takes the place of the description, tags, details
		// and status.
		return giode.Stack(st,
			header,
			fixed(descHeight+22+metaHeight+statusHeight+3*12, cd.chooser()),
			giode.Box(giode.Styles{FlexGrow: 1}),
			giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 6, Wrap: giode.WrapNoWrap},
				cd.download, cd.back),
		)
	}
	return giode.Stack(st,
		header,
		fixed(descHeight, giode.Text(tr.T(m.Description), giode.Styles{Color: "theme.secondary", FontSize: 13, LineHeight: 19, MaxLines: 3, TextOverflow: giode.TextOverflowEllipsis})),
		giode.Row(giode.Styles{Width: -1, Height: 22, Gap: 6, Wrap: giode.WrapNoWrap, Overflow: giode.OverflowHidden}, tags...),
		fixed(metaHeight, giode.Text(metaLine(m), giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 16, MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis})),
		fixed(statusHeight, cd.status()),
		giode.Box(giode.Styles{FlexGrow: 1}),
		cd.actions(),
	)
}
