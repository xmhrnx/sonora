// Command makeicon draws Sonora's app icon: an S for Sonora made of
// three parallel sound-wave strands, the top bowl turned one way and the
// bottom the other - speech going out and coming in - glowing violet to
// cyan on a deep navy tile.
//
//	go run ./tools/makeicon assets/icon
//
// writes icon.png (full bleed, for Windows, Linux and the window),
// icon-macos.png (inset with a shadow, on macOS's icon grid) and the
// layers of Android's adaptive icon: android-foreground.png (the S on a
// transparent 108dp canvas, inside the middle 66dp) and
// android-monochrome.png (its shape, tinted for themed icons). The
// background layer is the navy color set in giode.toml.
package main

import (
	"image/color"
	"os"
	"path/filepath"
)

const (
	size = 1024
	ss   = 3 // supersampling factor for smooth edges
)

var (
	violet = color.NRGBA{0x8b, 0x5c, 0xf6, 0xff}
	cyan   = color.NRGBA{0x22, 0xd3, 0xee, 0xff}
	white  = color.NRGBA{0xff, 0xff, 0xff, 0xff}
	// The tile: deep navy, a little lighter at the top.
	navyTop    = color.NRGBA{0x1e, 0x1b, 0x4b, 0xff}
	navyBottom = color.NRGBA{0x0b, 0x10, 0x26, 0xff}
)

func main() {
	out := "assets/icon"
	if len(os.Args) > 1 {
		out = os.Args[1]
	}
	save(filepath.Join(out, "icon.png"), render(tile{margin: 0}))
	// macOS icons sit on an 824px body inside the 1024px canvas.
	save(filepath.Join(out, "icon-macos.png"), render(tile{margin: 100, shadow: true}))
	// Android shows the middle 66 of the layer's 108dp: the S keeps to
	// it, and the launcher draws the navy background and the shape.
	save(filepath.Join(out, "android-foreground.png"), render(tile{none: true, scale: 66.0 / 108}))
	save(filepath.Join(out, "android-monochrome.png"), render(tile{none: true, scale: 66.0 / 108, mono: true}))
}
