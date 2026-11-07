package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"golang.org/x/image/draw"
)

// render draws the icon.
func render(t tile) image.Image {
	n := size * ss
	big := image.NewNRGBA(image.Rect(0, 0, n, n))
	m := float64(t.margin * ss)
	side := float64(n) - 2*m
	radius := side * 0.225
	// The S's unit square: the tile, or the middle of an empty canvas.
	art, ax := side, m
	if t.none {
		art = float64(n) * t.scale
		ax = (float64(n) - art) / 2
	}
	ss := strands()
	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			fx, fy := float64(x)+0.5, float64(y)+0.5
			var c color.NRGBA
			if !t.none {
				d := distToRoundRect(fx-m, fy-m, side, side, radius)
				if t.shadow && d > 0 {
					// A soft drop shadow below the tile.
					ds := distToRoundRect(fx-m, fy-m-side*0.02, side, side, radius)
					if ds > 0 && ds < side*0.05 {
						a := 0.30 * (1 - ds/(side*0.05))
						big.SetNRGBA(x, y, color.NRGBA{0, 0, 0, uint8(255 * a * a)})
					}
					continue
				}
				if d > 0 {
					continue
				}
				c = mix(navyTop, navyBottom, (fy-m)/side)
			}
			// The S, in the art square's unit coordinates.
			u, v := (fx-ax)/art, (fy-ax)/art
			dist, along := math.Inf(1), 0.0
			for _, s := range ss {
				if d, a := s.dist(u, v); d < dist {
					dist, along = d, a
				}
			}
			px := 1 / art // one pixel in unit coordinates, for antialiasing
			cover := clamp(0.5-dist/px, 0, 1)
			ink := mix(violet, cyan, along)
			// A highlight down the middle of each stroke.
			if dist < 0 {
				ink = mix(ink, white, 0.22*clamp(1+dist/(ss[0].width/2), 0, 1))
			}
			switch {
			case t.mono:
				c = color.NRGBA{255, 255, 255, uint8(255 * cover)}
			case t.none:
				c = ink
				c.A = uint8(255 * cover)
			default:
				// A faint glow around the strokes, then the strokes.
				if glow := clamp(1-dist/0.06, 0, 1); glow > 0 && dist > 0 {
					c = mix(c, ink, 0.22*glow*glow)
				}
				c = mix(c, ink, cover)
			}
			big.SetNRGBA(x, y, c)
		}
	}
	small := image.NewNRGBA(image.Rect(0, 0, size, size))
	draw.CatmullRom.Scale(small, small.Bounds(), big, big.Bounds(), draw.Src, nil)
	return small
}

func clamp(x, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, x)) }

// distToRoundRect is the signed distance from (x, y) to a rounded
// rectangle at the origin: negative inside.
func distToRoundRect(x, y, w, h, r float64) float64 {
	cx, cy := math.Abs(x-w/2)-(w/2-r), math.Abs(y-h/2)-(h/2-r)
	ox, oy := math.Max(cx, 0), math.Max(cy, 0)
	return math.Hypot(ox, oy) + math.Min(math.Max(cx, cy), 0) - r
}

func mix(a, b color.NRGBA, t float64) color.NRGBA {
	t = clamp(t, 0, 1)
	l := func(x, y uint8) uint8 { return uint8(float64(x) + (float64(y)-float64(x))*t + 0.5) }
	return color.NRGBA{l(a.R, b.R), l(a.G, b.G), l(a.B, b.B), l(a.A, b.A)}
}

func save(path string, img image.Image) {
	f, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}
