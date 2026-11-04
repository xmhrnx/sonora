package studio

import (
	"time"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/progressbar"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

func (v *voicePanel) footer() giode.Element {
	rt := v.rt
	c := v.p.c
	busy := c.Manager.Job != nil
	ready := rt.State == manager.StateReady || rt.State == manager.StateWorking
	st := giode.Styles{Width: -1, Justify: giode.JustifyCenter, BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 15, BorderRadius: theme.RadiusMd,
		Padding: giode.SymmetricInset(16, 13), BoxShadow: giode.BoxShadow{Y: 8, Blur: 24, Color: "#8b5cf650"}}
	if busy {
		st.BackgroundGradient = nil
		st.Background = "theme.danger"
		st.BoxShadow = giode.BoxShadow{}
		v.generate.Label(tr.T("Stop")).Icon(icons.Solid(icons.Stop, 14, ""))
	} else {
		v.generate.Label(tr.T("Generate speech")).Icon(icons.Icon(icons.Sparkles, 18, ""))
	}
	v.generate.Styles(st).Disabled(!busy && (!ready || v.p.editor.text() == ""))

	var status giode.Element
	switch {
	case busy:
		job := c.Manager.Job
		line := giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Spinner(giode.Styles{Width: 14, Height: 14, Color: "theme.primary"}),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T(job.Runtime.Status),
				giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis})),
			giode.Text(format.Clock(time.Since(job.Started)), giode.Styles{Color: "theme.muted", FontSize: 12, FontFamily: theme.Mono}),
		)
		if job.Runtime.Progress >= 0 {
			status = giode.Stack(giode.Styles{Width: -1, Gap: 8}, progressbar.New(job.Runtime.Progress), line)
		} else {
			status = line
		}
	case rt.Err != "" && rt.State == manager.StateReady:
		status = giode.Text(tr.T(format.FirstLine(rt.Err)), giode.Styles{Color: "theme.danger", FontSize: 12, MaxLines: 3})
	case !ready:
		status = giode.Text(tr.T("Install the engine to start generating."), theme.Caption)
	default:
		status = giode.Text(tr.T("Audio is saved to your Library automatically."), theme.Caption)
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10, Padding: giode.UniformInset(20)},
		giode.Box(giode.Styles{Width: -1}, v.generate),
		status,
	)
}
