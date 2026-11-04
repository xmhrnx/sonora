package profiles

import (
	"os"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/tr"
)

// importFile makes a profile from an audio file the user picks: cut to
// profiles.MaxLength, then transcribed by a ready transcription engine
// so cloning engines get the words too.
func (p *Page) importFile() {
	if p.importing {
		return
	}
	c := p.c
	p.importing = true
	var made profiles.Profile
	c.Background(func() error {
		f, err := c.PickFile(append([]string{"audio/*", "video/*"}, decode.Extensions...), "")
		if err != nil {
			return err
		}
		defer f.Close()
		src, err := core.SavePicked(f, core.PickedDir(), "voice")
		if err != nil {
			return err
		}
		defer os.Remove(src)
		made, err = c.Profiles.Add(title(f.Name), "", src, true)
		return err
	}, func(err error) {
		p.importing = false
		if err != nil {
			c.Error(err)
			return
		}
		if c.CanFillTranscripts() {
			c.Info(tr.T("Imported {name}; finding out what is said…", "name", made.Name))
			c.FillTranscript(made, nil)
			return
		}
		c.Success(tr.T("Imported {name}", "name", made.Name))
	})
}
