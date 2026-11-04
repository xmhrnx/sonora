package core

import (
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
)

// giode.toml's context actions, which installers register, must offer
// the same files as the app's.
func TestContextActionsMatchGiodeToml(t *testing.T) {
	data, err := os.ReadFile("../../giode.toml")
	if err != nil {
		t.Fatal(err)
	}
	block := regexp.MustCompile(`(?s)\[\[context_actions\]\]\s*id = "([a-z]+)".*?extensions = \[([^\]]*)\]`)
	found := map[string][]string{}
	for _, m := range block.FindAllStringSubmatch(string(data), -1) {
		var exts []string
		for _, e := range strings.Split(m[2], ",") {
			exts = append(exts, strings.Trim(strings.TrimSpace(e), `"`))
		}
		found[m[1]] = exts
	}
	for _, a := range ContextActions() {
		if !slices.Equal(found[a.ID], a.Extensions) {
			t.Errorf("giode.toml's %q action lists %v, the app %v", a.ID, found[a.ID], a.Extensions)
		}
	}
}
