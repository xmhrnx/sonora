package studio

import (
	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
)

// voiceChoice is the Voice section: the voice picker, plus a switch to
// cloning a voice profile on engines that can copy a voice.
func voiceChoice(cl *cloner, f *engineForm) giode.Element {
	m := f.rt.Manifest()
	switch m.Clone {
	case engines.CloneNone:
		return f.voiceTrigger.Element()
	case engines.CloneRequired:
		return cl.element(m)
	}
	body := f.voiceTrigger.Element()
	if f.cloneMode {
		body = cl.element(m)
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 12}, f.modes.Element(), body)
}
