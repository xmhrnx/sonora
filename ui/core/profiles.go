package core

import (
	"slices"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/profiles"
	"github.com/mehran9675/sonora/ui/tr"
)

// transcriber returns the speech-to-text engine that fills in
// profiles' transcripts: the one chosen on the Transcribe page if it's
// ready, else any ready one, else nil.
func (c *Ctx) transcriber() *manager.Runtime {
	ready := c.Manager.ReadyFor(engines.TaskTranscribe)
	if i := slices.IndexFunc(ready, func(rt *manager.Runtime) bool {
		return rt.Manifest().ID == c.Settings.Get().TranscribeEngine
	}); i >= 0 {
		return ready[i]
	}
	if len(ready) > 0 {
		return ready[0]
	}
	return nil
}

// CanFillTranscripts reports whether a transcription engine is ready to
// hear what's said in imported profiles.
func (c *Ctx) CanFillTranscripts() bool { return c.transcriber() != nil }

// FillTranscript has a transcription engine write down what is said in
// p's recording and stores it; done (optional) runs afterwards on the
// UI goroutine. Without a ready engine it does nothing.
func (c *Ctx) FillTranscript(p profiles.Profile, done func(error)) {
	rt := c.transcriber()
	if rt == nil {
		return
	}
	c.Manager.TranscribeClip(rt, c.Profiles.Path(p), func(text string, err error) {
		if err == nil {
			err = c.Profiles.Update(p.ID, func(p *profiles.Profile) { p.Transcript = text })
		}
		if err != nil {
			c.Error(err)
		} else {
			c.Success(tr.T("Transcript of {name} is ready", "name", p.Name))
		}
		if done != nil {
			done(err)
		}
	})
}
