package onnx

import (
	"context"
	"fmt"
	"strings"

	"github.com/mehran9675/sonora/internal/audio/wav"
	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// Transcribe cuts the recording into utterances where there is speech
// and recognizes each, so segments carry their times.
func (e *Recognizer) Transcribe(ctx context.Context, r engines.TranscribeRequest, log engines.Logger) (*engines.Transcript, error) {
	clip, err := wav.Read(r.Audio)
	if err != nil {
		return nil, err
	}
	if clip.SampleRate != 16000 {
		clip.Samples, clip.SampleRate = clip.Resample(16000), 16000
	}
	rec, vad, err := e.load(ctx, r, log)
	if err != nil {
		return nil, err
	}
	total := len(clip.Samples)
	out := &engines.Transcript{}
	last := -1
	var failure error
	recognize := func(start int, speech []float32) bool {
		if ctx.Err() != nil {
			return false
		}
		if pct := int(int64(start) * 100 / int64(max(total, 1))); pct != last {
			last = pct
			log(fmt.Sprintf("Transcribing… %d%%", pct))
		}
		res, err := rec.Recognize(speech, 16000)
		if err != nil {
			failure = err
			return false
		}
		if text := strings.TrimSpace(res.Text); text != "" {
			out.Segments = append(out.Segments, engines.Segment{
				Start: float64(start) / 16000, End: float64(start+len(speech)) / 16000, Text: text})
		}
		// Only Whisper's language detection is reliable enough to show.
		if out.Language == "" && e.pick(r.Options).Family == sherpa.Whisper {
			out.Language = strings.Trim(res.Language, "<|>")
		}
		return true
	}
	// The detector reports speech a little after it starts; widen each
	// stretch so first syllables aren't lost.
	const before, after = 16000 * 4 / 10, 16000 * 2 / 10
	vad.Split(clip.Samples, func(start int, speech []float32) bool {
		from, to := max(0, start-before), min(total, start+len(speech)+after)
		return recognize(from, clip.Samples[from:to])
	})
	// Speech too quiet for the detector: fall back to fixed windows.
	if len(out.Segments) == 0 && failure == nil && ctx.Err() == nil {
		step := maxUtterance * 16000
		for i := 0; i < total && recognize(i, clip.Samples[i:min(total, i+step)]); i += step {
		}
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if failure != nil {
		return nil, failure
	}
	log("Transcribing… 100%")
	parts := make([]string, len(out.Segments))
	for i, s := range out.Segments {
		parts[i] = s.Text
	}
	out.Text = JoinText(parts)
	return out, nil
}
