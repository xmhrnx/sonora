package transcribe

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/group"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/components/voicerecorder"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// panel is the right-hand column: engine, audio, language, settings
// and the Transcribe button.
type panel struct {
	p        *Page
	scroll   *scroll.Scroll
	engines  *picker.Picker
	engine   *picker.Trigger
	install  *giode.ButtonElement
	manage   *giode.ButtonElement
	pick     *giode.ButtonElement
	clear    *giode.ButtonElement
	play     *giode.PressableElement
	recorder *voicerecorder.Recorder
	run      *giode.ButtonElement
	eject    *giode.ButtonElement
	forms    map[string]*form
	rt       *manager.Runtime // engine shown this frame
}

func (v *panel) element(rt *manager.Runtime) giode.Element {
	v.rt = rt
	card := theme.Card()
	card.Width, card.Height, card.Gap = 340, -1, 0
	card.Padding = giode.UniformInset(0)
	return giode.Stack(card,
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1}, v.scroll),
		giode.Divider(giode.Styles{Color: "theme.border"}),
		v.footer(),
	)
}

// compactElement is the panel for the phone layout, without inner
// scrolling.
func (v *panel) compactElement(rt *manager.Runtime) giode.Element {
	v.rt = rt
	card := theme.Card()
	card.Width, card.Gap = -1, 0
	card.Padding = giode.UniformInset(0)
	return giode.Stack(card, v.content(), giode.Divider(giode.Styles{Color: "theme.border"}), v.footer())
}

func (v *panel) content() giode.Element {
	rt := v.rt
	els := []giode.Element{group.New(tr.T("Engine"), v.engineSection())}
	if callout := v.callout(); callout != nil {
		els = append(els, callout)
	}
	els = append(els, group.New(tr.T("Audio"), v.sourceSection()))
	if rt.State == manager.StateReady || rt.State == manager.StateWorking {
		f := v.form(rt)
		if f.language != nil {
			els = append(els, f.language.Element())
		}
		if f.translate != nil {
			els = append(els, f.translate.Element())
		}
		for _, fl := range f.fields {
			els = append(els, fl.Element())
		}
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 22, Padding: giode.UniformInset(20)}, els...)
}
