package profiledialog

import (
	"os"
	"strings"

	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/tr"
)

// Open shows the dialog with an empty name and the next passage.
func (d *Dialog) Open() {
	d.discard()
	d.name.SetText("")
	d.passage = (d.passage + 1) % len(profiles.Passages)
	d.open = true
}

// Close discards an unsaved take and closes the dialog.
func (d *Dialog) Close() {
	if d.recorder.Recording() {
		return // Stop or Cancel the recording first
	}
	d.discard()
	d.open = false
}

// discard deletes the take, if any.
func (d *Dialog) discard() {
	if d.take != "" {
		os.Remove(d.take)
		d.take = ""
	}
}

// profileName is the name typed, or a default.
func (d *Dialog) profileName() string {
	if n := strings.TrimSpace(d.name.Text()); n != "" {
		return n
	}
	return tr.T("My voice")
}

// saveProfile stores the take with the passage as its transcript.
func (d *Dialog) saveProfile() {
	if d.take == "" || d.saving {
		return
	}
	d.saving = true
	name, passage, take := d.profileName(), profiles.Passages[d.passage], d.take
	var p profiles.Profile
	d.c.Background(func() error {
		var err error
		p, err = d.c.Profiles.Add(name, passage, take, false)
		return err
	}, func(err error) {
		d.saving = false
		if err != nil {
			d.c.Error(err)
			return
		}
		d.discard()
		d.open = false
		d.c.Success(tr.T("Saved the profile {name}", "name", p.Name))
		if d.onSaved != nil {
			d.onSaved(p)
		}
	})
}
