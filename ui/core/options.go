package core

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/ui/components/optionfield"
	"github.com/mehran9675/sonora/ui/format"

	"github.com/mehran9675/sonora/ui/tr"
)

// OptionStore edits engine man's settings for its option fields. Values
// read as the manager resolves them, so an untouched model setting
// shows a model that's downloaded, and model choices say whether
// they're downloaded. onSet (optional) runs after a change.
func (c *Ctx) OptionStore(man *engines.Manifest, onSet func(id string)) optionfield.Store {
	return optionfield.Store{
		Get: func(id string) any { return c.Manager.Values(man)[id] },
		Set: func(id string, v any) {
			c.Settings.SetValue(man.ID, id, v)
			if onSet != nil {
				onSet(id)
			}
		},
		Detail: func(id, value string) string {
			o := man.ModelOption()
			if o == nil || o.ID != id {
				return ""
			}
			if c.Manager.Downloaded(man.ID)[value] {
				return tr.T("Downloaded")
			}
			return tr.T("Not downloaded yet")
		},
	}
}

// InstallModel installs an engine if needed and downloads one of its
// models, keeping the device awake meanwhile; done (optional) runs
// afterwards on the UI goroutine.
func (c *Ctx) InstallModel(rt *manager.Runtime, model string, done func(error)) {
	awake := c.StayAwake(tr.T("Downloading a model for {name}", "name", rt.Manifest().Name))
	c.Manager.InstallModel(rt, model, func(err error) {
		awake()
		if err == nil {
			c.NotifyIfAway(tr.T("{name} is ready", "name", rt.Manifest().Name), tr.T("The model is downloaded."))
		}
		if done != nil {
			done(err)
		}
	})
}

// SetupModel is what an Install button outside the Engines page
// downloads for rt: the model the page's selector shows, or the
// engine's only model. ok is false when the engine has several and the
// user picks one on the Engines page (ChooseModel); value is "" for an
// engine with no models to choose (a Python engine installs its
// environment).
func (c *Ctx) SetupModel(rt *manager.Runtime) (value, label string, ok bool) {
	m := rt.Manifest()
	models := engines.Models(rt.Engine)
	if o := m.ModelOption(); o != nil {
		value = c.Manager.Values(m).String(o.ID)
	} else if len(models) == 1 {
		value = models[0].Value
	} else if len(models) > 1 {
		return "", "", false
	}
	for _, ch := range models {
		if ch.Value == value {
			label = format.Parts(ch.Label)
		}
	}
	return value, label, true
}

// Setup installs rt from a page other than Engines: downloading the
// model SetupModel names, installing an engine without models, or
// opening the Engines page to choose one.
func (c *Ctx) Setup(rt *manager.Runtime, done func(error)) {
	value, _, ok := c.SetupModel(rt)
	switch {
	case !ok:
		if c.ChooseModel != nil {
			c.ChooseModel(rt.Manifest().ID)
		}
	case value == "":
		c.Install(rt, done)
	default:
		c.InstallModel(rt, value, done)
	}
}

// SetupLabel is the text of a Setup button: "Download Small · 466 MB",
// "Choose a model", or "Install".
func (c *Ctx) SetupLabel(rt *manager.Runtime) string {
	value, label, ok := c.SetupModel(rt)
	switch {
	case !ok:
		return tr.T("Choose a model")
	case value == "":
		return tr.T("Install")
	}
	return tr.T("Download {model}", "model", label)
}
