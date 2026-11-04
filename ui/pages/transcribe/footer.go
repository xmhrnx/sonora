package transcribe

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

func (v *panel) footer() giode.Element {
	rt := v.rt
	c := v.p.c
	job := c.Manager.TranscribeJob
	ready := rt.State == manager.StateReady || rt.State == manager.StateWorking
	st := giode.Styles{Width: -1, Justify: giode.JustifyCenter, BackgroundGradient: theme.Brand(90), Color: "#ffffff", FontSize: 15,
		BorderRadius: theme.RadiusMd, Padding: giode.SymmetricInset(16, 13), BoxShadow: giode.BoxShadow{Y: 8, Blur: 24, Color: "#8b5cf650"}}
	if job != nil {
		st.BackgroundGradient, st.Background, st.BoxShadow = nil, "theme.danger", giode.BoxShadow{}
		v.run.Label(tr.T("Stop")).Icon(icons.Solid(icons.Stop, 14, ""))
	} else {
		v.run.Label(tr.T("Transcribe")).Icon(icons.Icon(icons.Captions, 18, ""))
	}
	v.run.Styles(st).Disabled(job == nil && (!ready || v.p.source == "" || v.recorder.Recording()))

	var status giode.Element
	switch {
	case job != nil:
		line := giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Spinner(giode.Styles{Width: 14, Height: 14, Color: "theme.primary"}),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T(job.Runtime.Status),
				giode.Styles{Color: "theme.muted", FontSize: 12, TextOverflow: giode.TextOverflowEllipsis})),
			giode.Text(format.Clock(time.Since(job.Started)), giode.Styles{Color: "theme.muted", FontSize: 12, FontFamily: theme.Mono}),
		)
		status = line
		if job.Runtime.Progress >= 0 {
			status = giode.Stack(giode.Styles{Width: -1, Gap: 8}, progressbar.New(job.Runtime.Progress), line)
		}
	case rt.Err != "" && rt.State == manager.StateReady:
		status = giode.Text(tr.T(format.FirstLine(rt.Err)), giode.Styles{Color: "theme.danger", FontSize: 12, MaxLines: 3})
	case !ready:
		status = giode.Text(tr.T("Install the engine to start transcribing."), theme.Caption)
	default:
		status = giode.Text(tr.T("Transcripts are kept below, with their audio."), theme.Caption)
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10, Padding: giode.UniformInset(20)},
		giode.Box(giode.Styles{Width: -1}, v.run),
		status,
	)
}
