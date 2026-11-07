package transcribe

import (
	"os"

	"github.com/mehran9675/sonora/internal/manager"
	"github.com/mehran9675/sonora/internal/transcripts"
	"github.com/mehran9675/sonora/ui/tr"
)

// setSource replaces the audio waiting to be transcribed.
func (p *Page) setSource(path, name string, temp bool) {
	p.dropSource()
	p.source, p.sourceName, p.temp = path, name, temp
}

// dropSource forgets the waiting audio, deleting Sonora's own copy.
func (p *Page) dropSource() {
	if p.temp && p.source != "" {
		os.Remove(p.source)
	}
	p.source, p.sourceName, p.temp = "", "", false
}

// start transcribes the waiting audio, or stops the running job.
func (p *Page) start() {
	c := p.c
	if j := c.Manager.TranscribeJob; j != nil {
		c.Manager.Cancel(j.Runtime)
		return
	}
	rt := p.engine()
	if rt == nil || p.source == "" {
		return
	}
	s := c.Settings.Get()
	in := manager.TranscribeInput{Source: p.source, Title: p.sourceName,
		Language: language(rt.Manifest(), s.Language), Translate: s.Translate}
	awake := c.StayAwake(tr.T("Transcribing"))
	c.Manager.Transcribe(rt, in, func(e transcripts.Entry, err error) {
		awake()
		if err != nil {
			c.Error(err)
			return
		}
		c.NotifyIfAway(tr.T("Transcript ready"), e.Title)
		p.dropSource()
		p.result.show(e)
	})
}
