package settings

import (
	"github.com/mehran9675/sonora/ui/tr"
)

// limitLabel names a recording length in minutes.
func limitLabel(m int) string {
	switch {
	case m == 60:
		return tr.T("1 hour")
	case m%60 == 0:
		return tr.N("{n} hours", m/60)
	}
	return tr.N("{n} minutes", m)
}
