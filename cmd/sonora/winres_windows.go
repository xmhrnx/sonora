//go:build !giode_gio

package main

// Sonora's own Windows resources, for builds without the giode CLI (the
// CLI adds its own; see package winres).
import _ "github.com/mehran9675/sonora/cmd/sonora/winres"
