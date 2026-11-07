package transcribe

import (
	"time"

	"github.com/mehran9675/sonora/internal/transcripts"
)

// show selects e.
func (r *result) show(e transcripts.Entry) {
	r.flush()
	r.entry = &e
	r.text.SetText(e.Text)
}

// clear deselects the transcript (after it's deleted).
func (r *result) clear() {
	if r.save != nil {
		r.save.Stop()
	}
	r.entry = nil
	r.text.SetText("")
}

// edited saves the user's changes shortly after they stop typing.
func (r *result) edited(text string) {
	if r.entry == nil {
		return
	}
	r.entry.Text = text
	if r.save != nil {
		r.save.Stop()
	}
	id := r.entry.ID
	r.save = time.AfterFunc(600*time.Millisecond, func() { r.p.c.Transcripts.SetText(id, text) })
}

// flush saves pending edits now.
func (r *result) flush() {
	if r.save != nil && r.save.Stop() && r.entry != nil {
		r.p.c.Transcripts.SetText(r.entry.ID, r.entry.Text)
	}
}
