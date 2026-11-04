package studio

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/chip"
	"github.com/mehran9675/sonora/ui/components/progressbar"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// callout explains engines that can't generate yet.
func (v *voicePanel) callout() giode.Element {
	rt := v.rt
	m := rt.Manifest()
	var body []giode.Element
	switch rt.State {
	case manager.StateNotInstalled:
		msg := tr.T("{name} needs a one-time setup. Choose what to download; nothing downloads until you do.", "name", m.Name)
		body = append(body, giode.Text(msg, giode.Styles{Color: "theme.text", FontSize: 13, LineHeight: 19}))
		if m.Size != "" {
			body = append(body, chip.Chip(m.Size))
		}
		if fl := v.p.form(rt).modelField(); fl != nil {
			body = append(body, fl.Element())
		}
		if rt.Err != "" {
			body = append(body, giode.Text(tr.T(format.FirstLine(rt.Err)), giode.Styles{Color: "theme.danger", FontSize: 12, MaxLines: 3}))
		}
		v.install.Label(v.p.c.SetupLabel(rt))
		body = append(body, giode.Row(giode.Styles{Gap: 8}, v.install, v.manage))
	case manager.StateInstalling:
		body = append(body,
			giode.Text(tr.T("Installing {name}…", "name", m.Name), giode.Styles{Color: "theme.text", FontSize: 13, FontWeight: giode.FontWeightBold}),
			progressbar.New(rt.Progress),
			giode.Text(tr.T(rt.Status), giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis}),
		)
	case manager.StateUnavailable:
		body = append(body, giode.Text(tr.T(m.Hint), giode.Styles{Color: "theme.text", FontSize: 13, LineHeight: 19}))
	default:
		return nil
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 12, Padding: giode.UniformInset(14),
		Background: "theme.primary/0.08", BorderRadius: theme.RadiusMd, BorderWidth: 1, BorderColor: "theme.primary/0.25"}, body...)
}
