package transcribe

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/settings"
	"github.com/mehran9675/sonora/ui/components/optionfield"
	"github.com/mehran9675/sonora/ui/tr"
)

// form holds one engine's language, translation and option fields.
type form struct {
	rt        *manager.Runtime
	language  *optionfield.Field
	translate *optionfield.Field
	fields    []*optionfield.Field
}

// form returns the cached fields for rt.
func (v *panel) form(rt *manager.Runtime) *form {
	m := rt.Manifest()
	f := v.forms[m.ID]
	if f != nil && f.rt == rt {
		return f
	}
	c := v.p.c
	f = &form{rt: rt}
	global := optionfield.Store{
		Get: func(id string) any {
			s := c.Settings.Get()
			if id == "translate" {
				return s.Translate
			}
			return language(m, s.Language)
		},
		Set: func(id string, val any) {
			c.Settings.Update(func(s *settings.Settings) {
				switch id {
				case "translate":
					s.Translate, _ = val.(bool)
				default:
					s.Language, _ = val.(string)
				}
			})
		},
	}
	var choices []engines.Choice
	for _, l := range engines.Languages {
		if language(m, l.Value) == l.Value {
			choices = append(choices, l)
		}
	}
	if len(choices) > 1 {
		f.language = optionfield.New(engines.Option{ID: "language", Label: tr.Msg("Spoken language"), Type: engines.OptionSelect,
			Choices: choices, Help: tr.Msg("Detecting works well; naming the language helps short or noisy recordings.")}, global)
	}
	if m.Translate {
		f.translate = optionfield.New(engines.Option{ID: "translate", Label: tr.Msg("Translate into English"), Type: engines.OptionBool,
			Help: tr.Msg("Write the text in English, whatever language is spoken.")}, global)
	}
	store := c.OptionStore(m, nil)
	for _, o := range m.Options {
		if o.Scope != engines.ScopeSetup {
			f.fields = append(f.fields, optionfield.New(o, store))
		}
	}
	v.forms[m.ID] = f
	return f
}
