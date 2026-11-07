package engines

import (
	"strings"
)

// shortLabel drops the size from a model label: "Tiny · 39M · 110 MB"
// reads "Tiny".
func shortLabel(label string) string {
	name, _, _ := strings.Cut(label, " · ")
	return name
}
