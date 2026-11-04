package engines

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
)

func (cd *card) actions() giode.Element {
	rt := cd.rt
	m := rt.Manifest()
	var left, right []giode.Element
	switch rt.State {
	case manager.StateNotInstalled:
		left = append(left, cd.install)
	case manager.StateInstalling:
		left = append(left, cd.cancel)
	case manager.StateReady, manager.StateWorking:
		left = append(left, cd.use)
	}
	if rt.Engine.Loaded() && !rt.Busy() {
		right = append(right, cd.eject)
	}

	if m.NeedsInstall() || len(rt.Log) > 0 {
		right = append(right, cd.logBtn)
	}
	if m.Homepage != "" {
		right = append(right, cd.site)
	}
	if m.Custom {
		right = append(right, cd.configure)
	}
	installed := m.NeedsInstall() && (rt.State == manager.StateReady || rt.Err != "")
	if (installed || m.Custom) && !rt.Busy() {
		right = append(right, cd.remove)
	}
	return giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: 6, Wrap: giode.WrapNoWrap},
		append(append(left, giode.Box(giode.Styles{FlexGrow: 1})), right...)...)
}

func metaLine(m *engines.Manifest) string {
	parts := []string{}
	if m.Size != "" {
		parts = append(parts, m.Size)
	}
	if m.License != "" {
		parts = append(parts, m.License)
	}
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += "  ·  "
		}
		out += p
	}
	return out
}
