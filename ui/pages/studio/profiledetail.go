package studio

import (
	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/tr"
)

// profileDetail describes a profile in a list: its length and what is
// said.
func profileDetail(p profiles.Profile) string {
	said := tr.T("No transcript yet")
	if p.Transcript != "" {
		said = "“" + p.Transcript + "”"
	}
	return format.Join(format.Clock(p.Duration), said)
}
