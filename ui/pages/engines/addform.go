package engines

import (
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/ui/components/optionfield"
)

// setupForm is a set of setup fields editing a values map.
type setupForm struct {
	values engines.Values
	fields []*optionfield.Field
}

func newSetupForm(m *engines.Manifest) *setupForm {
	f := &setupForm{values: engines.Resolve(m.Options, nil)}
	st := optionfield.Store{Get: func(id string) any { return f.values[id] }, Set: func(id string, v any) { f.values[id] = v }}
	for _, o := range m.Options {
		if o.Scope == engines.ScopeSetup {
			f.fields = append(f.fields, optionfield.New(o, st))
		}
	}
	return f
}
