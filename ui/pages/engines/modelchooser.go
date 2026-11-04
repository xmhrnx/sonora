package engines

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/ui/tr"
)

// choose opens the model chooser, starting from the model the engine is
// set to use.
func (cd *card) choose() {
	m := cd.rt.Manifest()
	o := m.ModelOption()
	cd.choice = cd.c.Manager.Values(m).String(o.ID)
	// Adding a model: start from the first one that isn't here yet.
	if have := cd.c.Manager.Downloaded(m.ID); have[cd.choice] {
		for _, ch := range o.Choices {
			if !have[ch.Value] {
				cd.choice = ch.Value
				break
			}
		}
	}
	cd.choosing = true
}

// chooser asks which model to download.
func (cd *card) chooser() giode.Element {
	m := cd.rt.Manifest()
	msg := tr.T("You can add others later, or switch models in the Studio.")
	if m.Transcribes() {
		msg = tr.T("You can add others later, or switch models on the Transcribe page.")
	}
	cd.download.Disabled(cd.c.Manager.Downloaded(m.ID)[cd.choice])
	return giode.Stack(giode.Styles{Width: -1, Gap: 10},
		giode.Text(tr.T("Choose a model to download"), giode.Styles{Color: "theme.text", FontSize: 14, FontWeight: giode.FontWeightBold}),
		giode.Text(msg, giode.Styles{Color: "theme.secondary", FontSize: 12, LineHeight: 17, MaxLines: 3, TextOverflow: giode.TextOverflowEllipsis}),
		cd.model.Element(),
	)
}

// downloadedModels names the engine's downloaded models, or "" when it
// has no model choice.
func (cd *card) downloadedModels() string {
	o := cd.rt.Manifest().ModelOption()
	if o == nil {
		return ""
	}
	have := cd.c.Manager.Downloaded(cd.rt.Manifest().ID)
	var names []string
	for _, ch := range o.Choices {
		if have[ch.Value] {
			names = append(names, shortLabel(tr.T(ch.Label)))
		}
	}
	if len(names) == 0 {
		return ""
	}
	return tr.T("Downloaded: {models}", "models", strings.Join(names, ", "))
}

// allDownloaded reports whether every model of the engine is here.
func (cd *card) allDownloaded() bool {
	o := cd.rt.Manifest().ModelOption()
	return o != nil && len(cd.c.Manager.Downloaded(cd.rt.Manifest().ID)) >= len(o.Choices)
}

// choiceLabel returns the translated label of a model choice.
func choiceLabel(o *engines.Option, value string) string {
	for _, ch := range o.Choices {
		if ch.Value == value {
			return shortLabel(tr.T(ch.Label))
		}
	}
	return value
}

// shortLabel drops the size from a model label: "Tiny · 39M · 110 MB"
// reads "Tiny".
func shortLabel(label string) string {
	name, _, _ := strings.Cut(label, " · ")
	return name
}
