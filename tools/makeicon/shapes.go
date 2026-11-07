package main

import (
	"math"
)

// tile says how to draw: the tile inset by margin pixels (with a drop
// shadow), or no tile at all with the S scaled by scale, in color or as
// a white shape.
type tile struct {
	margin int
	shadow bool
	none   bool
	scale  float64
	mono   bool
}

// strand is one stroke of the S: two arcs of concentric ellipses, the
// top bowl and the bottom one, joined in the middle.
type strand struct {
	top, bottom arc
	width       float64
}

// arc is part of an ellipse centred at (cx, cy) with radii rx, ry, from
// angle a0 to a1 (radians, clockwise from 3 o'clock, y pointing down).
type arc struct {
	cx, cy, rx, ry, a0, a1 float64
}

// strands lays out the S in the unit square.
func strands() []strand {
	const (
		n     = 3
		width = 0.052 // stroke width
		step  = 0.083 // from one strand to the next
		inner = 0.050 // radius of the innermost strand
		sx    = 1.32  // the bowls are wider than tall
	)
	outer := inner + step*(n-1)
	// Strand i of the top bowl meets strand n-1-i of the bottom one in
	// the middle, so the bowls' centres are inner+outer apart.
	gap := inner + outer
	total := gap + 2*outer
	top := 0.5 - total/2 + outer
	bottom := top + gap
	var out []strand
	for i := 0; i < n; i++ {
		rt := inner + step*float64(i)
		rb := inner + step*float64(n-1-i)
		out = append(out, strand{
			// The top bowl runs from its foot in the middle, round the
			// left and over the top, to the upper right.
			top: arc{0.5, top, rt * sx, rt, deg(90), deg(330)},
			// The bottom bowl runs from its head in the middle, round the
			// right and under the bottom, to the lower left.
			bottom: arc{0.5, bottom, rb * sx, rb, deg(-90), deg(150)},
			width:  width,
		})
	}
	return out
}

func deg(d float64) float64 { return d * math.Pi / 180 }

// dist is the signed distance from (u, v) to the strand's edge
// (negative inside), and how far along the S the nearest point is, 0
// at its top end to 1 at its bottom end.
func (s strand) dist(u, v float64) (float64, float64) {
	dt, at := s.top.dist(u, v)
	db, ab := s.bottom.dist(u, v)
	if dt < db {
		// The top bowl is drawn from the middle backwards: its far end
		// is the S's start.
		return dt - s.width/2, 0.5 * (1 - at)
	}
	return db - s.width/2, 0.5 + 0.5*ab
}

// dist is the distance from (u, v) to the arc's centre line and the
// fraction of the arc where the nearest point lies. In the ellipse's
// squashed (circular) coordinates the radial gap is |ρ - ry|; dividing
// by the gradient's length turns it into the distance across the
// ellipse, exact enough near the stroke.
func (a arc) dist(u, v float64) (float64, float64) {
	k := a.rx / a.ry
	x, y := (u-a.cx)/k, v-a.cy
	t := math.Atan2(y, x)
	for t < a.a0 {
		t += 2 * math.Pi
	}
	for t > a.a0+2*math.Pi {
		t -= 2 * math.Pi
	}
	if t <= a.a1 {
		grad := math.Hypot(math.Cos(t)/k, math.Sin(t))
		return math.Abs(math.Hypot(x, y)-a.ry) / grad, (t - a.a0) / (a.a1 - a.a0)
	}
	// Past an end: the distance to the nearer end point (a round cap).
	end := func(t float64) float64 {
		return math.Hypot(u-(a.cx+a.rx*math.Cos(t)), v-(a.cy+a.ry*math.Sin(t)))
	}
	if d0, d1 := end(a.a0), end(a.a1); d0 < d1 {
		return d0, 0
	} else {
		return d1, 1
	}
}
