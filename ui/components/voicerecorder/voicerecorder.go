// Package voicerecorder records from the microphone - a voice to clone,
// or speech to transcribe: Record, a live meter and timer, Stop.
package voicerecorder

import (
	"time"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/record"
	"github.com/mehran9675/sonora/ui/components/levelmeter"
	"github.com/mehran9675/sonora/ui/core"
	"github.com/mehran9675/sonora/ui/icons"
	"github.com/mehran9675/sonora/ui/theme"
	"github.com/mehran9675/sonora/ui/tr"
)

// Options set what a recorder is for. The zero value records a voice
// for cloning.
type Options struct {
	// Min is the shortest take kept (default MinLength); Max stops the
	// recording (0: the limit set in Settings).
	Min, Max time.Duration
	Name     string // file name prefix; default "voice"
	// Hint shows while recording, Ready once the take is long enough;
	// both are English (tr.Msg) and shown translated.
	Hint, Ready string
}

// Recorder is stateful; create one per form.
type Recorder struct {
	c      *core.Ctx
	onDone func(path string)
	opts   Options

	rec     *record.Recorder
	opening bool // asking for the microphone or opening it
	timer   *time.Timer
	meter   *levelmeter.Meter
	start   *giode.ButtonElement
	stop    *giode.ButtonElement
	cancel  *giode.ButtonElement
}

const (
	// MaxLength stops a recording automatically.
	MaxLength = 30 * time.Second
	// MinLength is the shortest recording worth cloning from.
	MinLength = 5 * time.Second
)

// NewWith returns a recorder configured by opts.
func NewWith(c *core.Ctx, opts Options, onDone func(path string)) *Recorder {
	if opts.Min == 0 {
		opts.Min = MinLength
	}
	if opts.Name == "" {
		opts.Name = "voice"
	}
	if opts.Hint == "" {
		opts.Hint = tr.Msg("Read naturally, a few sentences, in a quiet room.")
	}
	if opts.Ready == "" {
		opts.Ready = tr.Msg("Sounds good — press Stop when you're done.")
	}
	r := &Recorder{c: c, onDone: onDone, opts: opts}
	r.meter = &levelmeter.Meter{Color: "theme.danger", Levels: func() []float32 {
		if r.rec == nil {
			return nil
		}
		return r.rec.Levels()
	}}
	r.start = giode.Button("", giode.Styles{Color: "theme.text", Background: "theme.border", BorderColor: "theme.border",
		FontSize: 13, Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}).
		Variant(giode.VariantOutlined).Icon(icons.Icon(icons.Mic, 15, "")).OnClick(r.begin)
	r.stop = giode.Button("", giode.Styles{Background: "theme.danger", Color: "#ffffff", FontSize: 13,
		Padding: giode.SymmetricInset(12, 7), BorderRadius: theme.RadiusSm}).
		Icon(icons.Solid(icons.Stop, 12, "")).OnClick(func() { r.finish(true) })
	tr.Bind(func(s string) { r.start.Label(s) }, "Record")
	tr.Bind(func(s string) { r.stop.Label(s) }, "Stop")
	r.cancel = giode.Button("", giode.Styles{Color: "theme.muted", Background: "theme.text", Padding: giode.UniformInset(6), BorderRadius: 8}).
		Variant(giode.VariantGhost).Icon(icons.Icon(icons.Close, 14, "")).OnClick(func() { r.finish(false) })
	return r
}
