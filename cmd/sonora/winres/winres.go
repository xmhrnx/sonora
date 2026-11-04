// Package winres carries Sonora's Windows resources - the icon, version
// information and a manifest with per-monitor DPI awareness, long paths
// and common controls v6 - as .syso files the linker picks up when
// this package is part of the build. cmd/sonora imports it only when
// built without the giode CLI, whose windows target writes resources of
// its own: two sets don't link ("too many .rsrc sections").
//
// Regenerate the .syso files after changing winres.json or the icon:
//
//	go install github.com/tc-hib/go-winres@v0.3.3
//	(cd cmd/sonora/winres && go-winres make --in winres.json --arch amd64,arm64)
package winres
