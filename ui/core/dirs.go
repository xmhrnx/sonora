package core

import "github.com/mehran9675/sonora/internal/paths"

// PickedDir holds copies of audio files picked on phones, kept until
// they're transcribed.
func PickedDir() string { return paths.Sub("cache", "picked") }

// CacheDir holds temporary files: shared exports, converted audio.
func CacheDir() string { return paths.Sub("cache") }
