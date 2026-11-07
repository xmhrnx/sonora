package voicerecorder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/record"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/format"
	"github.com/mehran9675/sonora/ui/platform"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// New returns a recorder for cloning that calls onDone with the saved
// file. Cloning models use a few seconds of speech, so it stops at
// MaxLength.
func New(c *core.Ctx, onDone func(path string)) *Recorder {
	return NewWith(c, Options{Max: MaxLength}, onDone)
}

// Recording reports whether the microphone is open.
func (r *Recorder) Recording() bool { return r.rec != nil }

// Button is the Record button shown while idle.
func (r *Recorder) Button() giode.Element { return r.start.Loading(r.opening) }

func (r *Recorder) begin() {
	if r.opening {
		return
	}
	r.opening = true
	r.c.Player.Pause() // don't record the speakers
	// Recordings go to the folder chosen in Settings, at its quality.
	path := filepath.Join(r.c.RecordingsDir(), r.opts.Name+"-"+time.Now().Format("2006-01-02 15.04.05")+".opus")
	limit := r.limit()
	var rec *record.Recorder
	r.c.Background(func() error {
		// Android asks the user for the microphone the first time.
		if err := platform.RequestMicrophone(r.c.App, context.Background()); err != nil {
			return err
		}
		var err error
		rec, err = record.Start(record.Options{Path: path, Bitrate: r.c.RecordBitrate(), Device: r.c.Settings.Get().Microphone})
		return err
	}, func(err error) {
		r.opening = false
		if err != nil {
			r.c.Error(err)
			return
		}
		r.rec = rec
		r.timer = time.AfterFunc(limit, func() { r.c.App.Post(func() { r.finish(true) }) })
	})
}

// limit returns how long this recorder records at most.
func (r *Recorder) limit() time.Duration {
	if r.opts.Max > 0 {
		return r.opts.Max
	}
	return r.c.RecordMax()
}

// finish closes the microphone and the file and, when keep is set,
// hands the take over; otherwise it's deleted.
func (r *Recorder) finish(keep bool) {
	rec := r.rec
	if rec == nil {
		return
	}
	r.rec = nil
	r.timer.Stop()
	var path string
	r.c.Background(func() error {
		if !keep {
			rec.Cancel()
			return nil
		}
		p, length, err := rec.Stop()
		if err != nil {
			return err
		}
		if length < r.opts.Min {
			os.Remove(p)
			return errors.New(tr.N("That recording is too short; speak for at least {n} seconds", int(r.opts.Min.Seconds())))
		}
		path = p
		return nil
	}, func(err error) {
		if err != nil {
			r.c.Error(err)
			return
		}
		if path != "" {
			r.onDone(path)
		}
	})
}

// Element is the live recording panel; lay it out while Recording.
func (r *Recorder) Element() giode.Element {
	if r.rec == nil {
		return giode.Box(giode.Styles{})
	}
	giode.HandleBack(func() { r.finish(false) }) // Back discards the take
	elapsed := r.rec.Elapsed()
	clock := fmt.Sprintf("%s / %s", format.Clock(elapsed), format.Clock(r.limit()))
	hint := tr.T(r.opts.Hint)
	if elapsed >= r.opts.Min {
		hint = tr.T(r.opts.Ready)
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 10, Padding: giode.UniformInset(12),
		Background: "theme.danger/0.08", BorderRadius: theme.RadiusMd, BorderWidth: 1, BorderColor: "theme.danger/0.3"},
		giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{Width: 9, Height: 9, BorderRadius: 99, Background: "theme.danger"}),
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(tr.T("Recording"), giode.Styles{Color: "theme.text", FontSize: 13, FontWeight: giode.FontWeightBold})),
			giode.Text(clock, giode.Styles{Color: "theme.secondary", FontSize: 12, FontFamily: theme.Mono}),
			r.cancel,
		),
		giode.Box(giode.Styles{Width: -1}, r.meter),
		giode.Row(giode.Styles{Width: -1, Gap: 8, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(hint, giode.Styles{Color: "theme.muted", FontSize: 12})),
			r.stop,
		),
	)
}
