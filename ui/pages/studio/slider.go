package studio

import (
	"github.com/xmhrnx/giode"
)

// syncSlider keeps the speed slider in step with settings changed
// elsewhere (reusing a library entry), except while it's dragged.
func syncSlider(f *engineForm) giode.Element {
	if !f.speed.Dragging() {
		f.speed.SetValue(speedToSlider(f.speedValue()))
	}
	return giode.Box(giode.Styles{Width: -1}, f.speed)
}
