// Package tr translates Sonora's text. Messages are keyed by their
// English text, so English needs no catalog except for plural forms:
//
//	tr.T("Recording")                         // "Aufnahme" in German
//	tr.T("Downloading {file}", "file", name)  // placeholders by name
//	tr.N("{n} words", count)                  // plural forms; {n} is count
//
// Text laid out every frame calls T where it's drawn. Elements built
// once (buttons, inputs) take their text through Bind, which sets it
// again whenever the language changes.
package tr

import (
	"embed"
)

//go:embed locales/*.json
var files embed.FS
