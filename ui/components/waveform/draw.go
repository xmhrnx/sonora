package waveform

import (
	"image"
	"image/color"

	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
)

func (w *Waveform) handleInput(gtx layout.Context, size image.Point, enabled bool) {
	for {
		ev, ok := gtx.Event(pointer.Filter{Target: w, Kinds: pointer.Press | pointer.Drag})
		if !ok {
			return
		}
		e, ok := ev.(pointer.Event)
		if !ok || !enabled || w.OnSeek == nil {
			continue
		}
		w.OnSeek(max(0, min(1, e.Position.X/float32(max(1, size.X)))))
		gtx.Execute(op.InvalidateCmd{})
	}
}

// bucket returns the loudest peak among those falling in bar i of n.
func bucket(peaks []float32, i, n int) float32 {
	lo := i * len(peaks) / n
	hi := max(lo+1, (i+1)*len(peaks)/n)
	var v float32
	for _, p := range peaks[lo:min(hi, len(peaks))] {
		v = max(v, p)
	}
	return 0.08 + 0.92*v
}

func mix(a, b color.NRGBA, t float32) color.NRGBA {
	l := func(x, y uint8) uint8 { return uint8(float32(x) + (float32(y)-float32(x))*t) }
	return color.NRGBA{R: l(a.R, b.R), G: l(a.G, b.G), B: l(a.B, b.B), A: l(a.A, b.A)}
}
