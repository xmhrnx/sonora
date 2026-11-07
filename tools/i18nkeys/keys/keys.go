// Package keys finds the text Sonora translates: the keys of tr.T,
// tr.N, tr.Msg and tr.Bind calls in the UI code, the status and hint
// messages the engine manager reports, and the descriptions, option
// labels and language names engines show.
package keys

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mehran9675/sonora/internal/catalog"
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx"
)

// Key is a message to translate.
type Key struct {
	Text   string `json:"text"`
	Plural bool   `json:"plural,omitempty"` // used with tr.N
	Where  string `json:"where"`            // the first place it's used
}

// Find returns the keys under the repository root, sorted.
func Find(root string) ([]Key, error) {
	found := map[string]*Key{}
	add := func(text, where string, plural bool) {
		if strings.TrimSpace(text) == "" {
			return
		}
		k := found[text]
		if k == nil {
			k = &Key{Text: text, Where: where}
			found[text] = k
		}
		k.Plural = k.Plural || plural
	}
	if err := goKeys(filepath.Join(root, "ui"), root, add); err != nil {
		return nil, err
	}
	if err := managerKeys(filepath.Join(root, "internal", "manager"), root, add); err != nil {
		return nil, err
	}
	manifests, _ := catalog.LoadDir(filepath.Join(root, "python", "engines"))
	tmp, err := os.MkdirTemp("", "i18nkeys")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmp)
	for _, e := range onnx.All(tmp) {
		manifests = append(manifests, e.Manifest())
	}
	for _, m := range manifests {
		manifestKeys(m, add)
	}
	for _, l := range engines.Languages {
		add(l.Label, "engines.Languages", false)
	}
	out := make([]Key, 0, len(found))
	for _, k := range found {
		out = append(out, *k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Text < out[j].Text })
	return out, nil
}
