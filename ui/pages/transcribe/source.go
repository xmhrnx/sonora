package transcribe

import (
	"os"
	"path/filepath"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/decode"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// pickFile asks for an audio (or video) file. Desktop pickers return
// the file itself; Android's hands over a stream, copied first.
func (v *panel) pickFile() {
	c := v.p.c
	var path, name string
	var temp bool
	c.Background(func() error {
		f, err := c.PickFile(append([]string{"audio/*", "video/*"}, decode.Extensions...), tr.Msg("Record instead."))
		if err != nil {
			return err
		}
		defer f.Close()
		// Desktop pickers hand over the file itself: read it in place.
		if file, ok := f.ReadCloser.(*os.File); ok && filepath.IsAbs(file.Name()) {
			path, name = file.Name(), f.Name
			return nil
		}
		name = f.Name
		if name == "" {
			name = tr.T("Audio from {time}", "time", time.Now().Format("Jan 2, 15:04"))
		}
		path, err = core.SavePicked(f, core.PickedDir(), "audio")
		temp = true
		return err
	}, func(err error) {
		if err != nil {
			c.Error(err)
			return
		}
		v.p.setSource(path, name, temp)
	})
}

// sourceSection picks or records the audio, or shows the chosen file.
func (v *panel) sourceSection() giode.Element {
	p := v.p
	switch {
	case v.recorder.Recording():
		return v.recorder.Element()
	case p.source == "":
		return giode.Stack(giode.Styles{Width: -1, Gap: 10},
			giode.Text(tr.T("Pick a recording or video, or record right here."), giode.Styles{Color: "theme.muted", FontSize: 12, LineHeight: 17}),
			giode.Row(giode.Styles{Width: -1, Gap: 8}, v.pick, v.recorder.Button()),
		)
	}
	return giode.Row(giode.Styles{Width: -1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
		Padding: giode.SymmetricInset(10, 8), Background: "theme.primary/0.1", BorderRadius: theme.RadiusSm + 2},
		v.play,
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(p.sourceName,
			giode.Styles{Color: "theme.text", FontSize: 13, TextOverflow: giode.TextOverflowEllipsis})),
		v.clear,
	)
}
