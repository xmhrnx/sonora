// Command fetchsherpa downloads the sherpa-onnx native libraries that
// the on-device engines link against into third_party/sherpa-onnx.
//
//	go run ./tools/fetchsherpa android     # all Android ABIs, and android.aar
//	go run ./tools/fetchsherpa linux-x64   # for developing on Linux
//
// Android builds link against them automatically, and giode packs
// third_party/sherpa-onnx/android.aar (giode.toml's android.libraries)
// into the APK and the Play bundle; desktop builds need -tags sherpa.
// See internal/engines/onnx/sherpa.
package main

// Version is the sherpa-onnx release the binding's header matches.
const Version = "1.13.8"

// libs are the files each target needs.
var libs = []string{"libsherpa-onnx-c-api.so", "libonnxruntime.so"}
