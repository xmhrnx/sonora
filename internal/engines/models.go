package engines

import "context"

// ModelStore is implemented by engines whose models download one by
// one, each a choice the user makes (see Models): the engine card asks
// which to download, and the main page defaults to one that is.
type ModelStore interface {
	// DownloadedModels reports which of the models are on this
	// device.
	DownloadedModels() map[string]bool
	// InstallModel installs the engine if it isn't yet, then downloads
	// the model with the given choice value.
	InstallModel(ctx context.Context, model string, log Logger) error
}

// ModelLister is implemented by engines that download packages which
// aren't choices of a model option - Piper's voices, an engine with a
// single model - so the engine card can still ask which to download.
type ModelLister interface {
	// Models lists the packages, value and label (with its size) each.
	Models() []Choice
}

// Models returns what the user picks from to download: the model
// option's choices, or the engine's own list. It is empty for engines
// with nothing to choose.
func Models(e Engine) []Choice {
	if o := e.Manifest().ModelOption(); o != nil {
		return o.Choices
	}
	if l, ok := e.(ModelLister); ok {
		return l.Models()
	}
	return nil
}

// ModelOption returns the option that picks the engine's model: the
// one named by model_option, or a select named "model". It returns nil
// when the model isn't a choice from a list.
func (m *Manifest) ModelOption() *Option {
	id := m.ModelOptionID
	if id == "" {
		id = "model"
	}
	if o := m.Option(id); o != nil && o.Type == OptionSelect && len(o.Choices) > 0 {
		return o
	}
	return nil
}
