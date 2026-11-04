package engines

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/ui/components/logdialog"
	"github.com/mehran9675/sonora/ui/components/picker"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// initChooser builds the card's model chooser: a picker of the
// engine's models (see engines.Models), Download and Cancel, and the
// Add model button shown once one is here. Nothing downloads until the
// user picks a model.
func (cd *card) initChooser(primary, outline giode.Styles, logs *logdialog.Dialog) {
	c, rt, m := cd.c, cd.rt, cd.rt.Manifest()
	cd.models = picker.New(tr.Msg("Choose a model"), func() []picker.Item {
		have := c.Manager.Downloaded(m.ID)
		var items []picker.Item
		for _, ch := range engines.Models(rt.Engine) {
			it := picker.Item{ID: ch.Value, Label: format.Parts(ch.Label)}
			if have[ch.Value] {
				it.Detail = tr.T("Downloaded")
			}
			items = append(items, it)
		}
		return items
	}, func() string { return cd.choice }, func(it picker.Item) { cd.choice = it.ID })
	cd.model = picker.NewTrigger(tr.Msg("Model"), func() (string, string) {
		if cd.choice == "" {
			return tr.T("Choose a model"), ""
		}
		detail := ""
		if c.Manager.Downloaded(m.ID)[cd.choice] {
			detail = tr.T("Already downloaded")
		}
		return cd.modelLabel(cd.choice), detail
	}, nil, cd.models.Open)
	cd.download = giode.Button("", primary).Icon(icons.Icon(icons.Download, 15, "")).OnClick(func() {
		cd.choosing = false
		c.InstallModel(rt, cd.choice, func(err error) {
			if err != nil {
				c.Error(err)
				logs.Show(rt)
				return
			}
			c.Success(tr.T("{engine} is ready with {model}", "engine", m.Name, "model", shortLabel(cd.modelLabel(cd.choice))))
		})
	})
	cd.back = giode.Button("", outline).Variant(giode.VariantOutlined).OnClick(func() { cd.choosing = false })
	cd.addModel = giode.Button("", giode.Styles{Color: "theme.primary", Background: "theme.primary", FontSize: 12,
		Padding: giode.SymmetricInset(8, 5), BorderRadius: theme.RadiusSm}).Variant(giode.VariantGhost).
		Icon(icons.Icon(icons.Plus, 13, "")).OnClick(cd.choose)
	tr.Bind(func(s string) { cd.download.Label(s) }, "Download")
	tr.Bind(func(s string) { cd.back.Label(s) }, "Cancel")
	tr.Bind(func(s string) { cd.addModel.Label(s) }, "Add model")
}

// choose opens the model chooser with nothing picked: the user chooses
// what to download.
func (cd *card) choose() {
	cd.choice = ""
	cd.choosing = true
}

// chooser asks which model to download.
func (cd *card) chooser() giode.Element {
	m := cd.rt.Manifest()
	msg := tr.T("You can add others later, or switch models in the Studio.")
	if m.Transcribes() {
		msg = tr.T("You can add others later, or switch models on the Transcribe page.")
	}
	cd.download.Disabled(cd.choice == "" || cd.c.Manager.Downloaded(m.ID)[cd.choice])
	return giode.Stack(giode.Styles{Width: -1, Gap: 10},
		giode.Text(tr.T("Choose a model to download"), giode.Styles{Color: "theme.text", FontSize: 14, FontWeight: giode.FontWeightBold}),
		giode.Text(msg, giode.Styles{Color: "theme.secondary", FontSize: 12, LineHeight: 17, MaxLines: 3, TextOverflow: giode.TextOverflowEllipsis}),
		cd.model.Element(),
	)
}

// downloadedModels names the engine's downloaded models, or "" when
// none is (or it has none to choose).
func (cd *card) downloadedModels() string {
	have := cd.c.Manager.Downloaded(cd.rt.Manifest().ID)
	var labels, names []string
	for _, ch := range engines.Models(cd.rt.Engine) {
		if have[ch.Value] {
			labels = append(labels, format.Parts(ch.Label))
			names = append(names, shortLabel(format.Parts(ch.Label)))
		}
	}
	if len(names) == 0 {
		return ""
	}
	// Models that share a name (Piper's voices in two qualities) add
	// the part that tells them apart: "Lessac (High quality)".
	seen := map[string]int{}
	for _, n := range names {
		seen[n]++
	}
	for i, n := range names {
		if seen[n] > 1 {
			parts := strings.Split(labels[i], " · ")
			if len(parts) > 2 {
				names[i] = n + " (" + parts[len(parts)-2] + ")"
			}
		}
	}
	return tr.T("Downloaded: {models}", "models", strings.Join(names, ", "))
}

// allDownloaded reports whether every model of the engine is here.
func (cd *card) allDownloaded() bool {
	models := engines.Models(cd.rt.Engine)
	return len(models) > 0 && len(cd.c.Manager.Downloaded(cd.rt.Manifest().ID)) >= len(models)
}

// modelLabel returns a model's translated label, size included.
func (cd *card) modelLabel(value string) string {
	for _, ch := range engines.Models(cd.rt.Engine) {
		if ch.Value == value {
			return format.Parts(ch.Label)
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
