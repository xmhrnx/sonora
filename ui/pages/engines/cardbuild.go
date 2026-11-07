package engines

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/ui/components/logdialog"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

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
