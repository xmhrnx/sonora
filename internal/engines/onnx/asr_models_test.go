package onnx

import "testing"

// sherpa-onnx ends the process on a language a model doesn't know, so
// only known hints may reach it.
func TestSupportedLanguage(t *testing.T) {
	tiny, turbo, small := whisperModels[0], whisperModels[3], whisperModels[2]
	sense, parakeet := senseVoiceModels[0], parakeetModels[0]
	cases := []struct {
		m          model
		lang, want string
	}{
		{tiny, "de", "de"}, {tiny, "yue", ""}, {small, "yue", ""}, {turbo, "yue", "yue"}, {tiny, "xx", ""}, {tiny, "", ""},
		{sense, "ja", "ja"}, {sense, "de", ""}, {parakeet, "de", ""}, {moonshineModels[0], "en", ""},
	}
	for _, c := range cases {
		if got := supportedLanguage(c.m, c.lang); got != c.want {
			t.Errorf("%s %q: got %q, want %q", c.m.ID, c.lang, got, c.want)
		}
	}
}
