package keys

import (
	"strings"

	"github.com/mehran9675/sonora/internal/engines"
)

// manifestKeys collects what an engine shows besides names: its
// description, hint and tags, its voices' groups and descriptions, its options' labels, help and placeholders,
// and the labels of choices other than voices and models.
func manifestKeys(m *engines.Manifest, add func(text, where string, plural bool)) {
	where := "engine " + m.ID
	add(m.Description, where, false)
	add(m.Hint, where, false)
	for _, t := range m.Tags {
		add(t, where, false)
	}
	for _, v := range m.Voices {
		// Shown through core.Parts: each " · " part translates alone.
		for _, s := range []string{v.Group, v.Description} {
			for _, part := range strings.Split(s, " · ") {
				if part != "" && (part[0] < '0' || part[0] > '9') {
					add(part, where, false)
				}
			}
		}
	}
	for _, o := range m.Options {
		add(o.Label, where, false)
		add(o.Help, where, false)
		add(o.Placeholder, where, false)
		if o.ID == "voice" || o.ID == "speaker" || o.ID == m.ModelOptionID || (m.ModelOptionID == "" && o.ID == "model") {
			continue
		}
		for _, c := range o.Choices {
			add(c.Label, where, false)
		}
	}
}
