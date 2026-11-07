package transcribe

import (
	"strings"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/ui/tr"
)

// languageName names a language code ("de" → "German"), translated.
func languageName(code string) string {
	for _, l := range engines.Languages {
		if l.Value != "" && l.Value == strings.ToLower(code) {
			return tr.T(l.Label)
		}
	}
	return strings.ToUpper(code)
}
