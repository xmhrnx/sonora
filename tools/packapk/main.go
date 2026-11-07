// Command packapk adds the sherpa-onnx native libraries to an APK built
// by the giode CLI (which packs only the app's own libgio.so), then
// aligns and signs it again with the Android SDK's zipalign and
// apksigner.
//
//	go run ./tools/fetchsherpa android      # once: the libraries
//	giode build android                     # dist/android/Sonora.apk
//	go run ./tools/packapk dist/android/Sonora.apk
//
// Without -ks it signs with the SDK's debug key (~/.android/debug.keystore),
// for installing on your own phone. Releases pass their keystore, with
// its password in the environment variable named by -pass-env.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	libs := flag.String("libs", filepath.Join("third_party", "sherpa-onnx", "android"), "folder of <abi>/*.so to add")
	out := flag.String("o", "", "output APK (default: replace the input)")
	ks := flag.String("ks", "", "keystore to sign with (default: the Android debug keystore)")
	passEnv := flag.String("pass-env", "", "environment variable holding the keystore password (default: the debug key's)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: packapk [flags] app.apk")
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	in := flag.Arg(0)
	if *out == "" {
		*out = in
	}
	tools, err := buildTools()
	if err != nil {
		log.Fatal(err)
	}
	unsigned := *out + ".unaligned"
	added, err := addLibs(in, unsigned, *libs)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(unsigned)
	log.Printf("added %d libraries", added)
	if err := run(filepath.Join(tools, "zipalign"), "-p", "-f", "4", unsigned, *out); err != nil {
		log.Fatal(err)
	}
	sign := []string{"sign"}
	if *ks == "" {
		home, _ := os.UserHomeDir()
		sign = append(sign, "--ks", filepath.Join(home, ".android", "debug.keystore"), "--ks-pass", "pass:android")
	} else {
		if *passEnv == "" {
			log.Fatal("-ks needs -pass-env")
		}
		sign = append(sign, "--ks", *ks, "--ks-pass", "env:"+*passEnv)
	}
	apksigner := filepath.Join(tools, "apksigner")
	if err := run(apksigner, append(sign, *out)...); err != nil {
		log.Fatal(err)
	}
	if err := run(apksigner, "verify", *out); err != nil {
		log.Fatal(err)
	}
	log.Printf("wrote %s", *out)
}
