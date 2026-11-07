package profiles

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// saveName stores the name typed and leaves editing.
func (cd *card) saveName() {
	if name := strings.TrimSpace(cd.name.Text()); name != "" {
		cd.c.Error(cd.c.Profiles.Update(cd.id, func(p *profiles.Profile) { p.Name = name }))
	}
	cd.editing = false
}

// transcript shows what is said in the profile, or how to find out.
func (cd *card) transcript(pr profiles.Profile) giode.Element {
	quote := giode.Styles{Color: "theme.secondary", FontSize: 13, LineHeight: 19, MaxLines: 2, TextOverflow: giode.TextOverflowEllipsis}
	switch {
	case pr.Transcript != "":
		return giode.Text("“"+pr.Transcript+"”", quote)
	case cd.filling:
		return giode.Text(tr.T("Finding out what is said…"), theme.Caption)
	case cd.c.CanFillTranscripts():
		return giode.Row(giode.Styles{Gap: 8, Align: giode.AlignCenter},
			giode.Text(tr.T("No transcript yet."), theme.Caption), cd.transcribe)
	}
	return giode.Text(tr.T("No transcript yet. Install a transcription engine, such as Whisper, and Sonora fills it in; "+
		"engines that need one can't clone this voice until then."), giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 17})
}
