package studio

import (
	"github.com/mehran9675/sonora/ui/tr"
)

func (f *engineForm) speedValue() float64 {
	if s, ok := f.c.Settings.Get().Speed[f.rt.Manifest().ID]; ok && s > 0 {
		return s
	}
	return 1
}

func sliderToSpeed(v float32) float64 { return float64(int((0.5+float64(v)*1.5)*20+0.5)) / 20 }

func speedToSlider(s float64) float32 { return float32((s - 0.5) / 1.5) }

func orDefault(s string) string {
	if s == "" {
		return tr.T("Default voice")
	}
	return s
}
