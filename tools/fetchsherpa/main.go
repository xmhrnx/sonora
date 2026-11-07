// Command fetchsherpa downloads the sherpa-onnx native libraries that
// the on-device engines link against into third_party/sherpa-onnx.
//
//	go run ./tools/fetchsherpa android     # all Android ABIs, for the APK
//	go run ./tools/fetchsherpa linux-x64   # for developing on Linux
//
// Android builds use them automatically (see tools/packapk); desktop
// builds need -tags sherpa. See internal/engines/onnx/sherpa.
package main

// Version is the sherpa-onnx release the binding's header matches.
const Version = "1.13.8"

// libs are the files each target needs.
var libs = []string{"libsherpa-onnx-c-api.so", "libonnxruntime.so"}
