package onnx

import (
	"slices"
	"strings"

	"github.com/mehran9675/sonora/internal/engines/onnx/sherpa"
)

// whisperTurbo and whisperLargeV3 also know Cantonese (see
// supportedLanguage).
const (
	whisperTurbo   = "sherpa-onnx-whisper-turbo"
	whisperLargeV3 = "sherpa-onnx-whisper-large-v3"
)

// The speech recognition packages. Sizes are downloads; Whisper
// packages carry fp32 and int8 models, and only int8 is unpacked.
var (
	whisperModels = []model{
		{ID: "sherpa-onnx-whisper-tiny", Label: "Tiny · 39M", Size: "110 MB"},
		{ID: "sherpa-onnx-whisper-base", Label: "Base · 74M", Size: "197 MB"},
		{ID: "sherpa-onnx-whisper-small", Label: "Small · 244M", Size: "609 MB"},
		{ID: whisperTurbo, Label: "Large v3 Turbo · 809M", Size: "537 MB"},
		{ID: "sherpa-onnx-whisper-medium", Label: "Medium · 769M", Size: "1.9 GB"},
		{ID: whisperLargeV3, Label: "Large v3 · 1550M", Size: "1.1 GB"},
	}
	// parakeetModels: v3 knows 25 European languages and finds which
	// is spoken; v2 is English only, a little more accurate there.
	parakeetModels = []model{
		{ID: "sherpa-onnx-nemo-parakeet-tdt-0.6b-v3-int8", Label: "TDT 0.6B v3 · 25 languages", Size: "464 MB", Family: sherpa.Transducer},
		{ID: "sherpa-onnx-nemo-parakeet-tdt-0.6b-v2-int8", Label: "TDT 0.6B v2 · English", Size: "482 MB", Family: sherpa.Transducer},
	}
	senseVoiceModels = []model{
		{ID: "sherpa-onnx-sense-voice-zh-en-ja-ko-yue-int8-2025-09-09", Label: "Small · int8", Size: "158 MB", Family: sherpa.SenseVoice},
		{ID: "sherpa-onnx-sense-voice-zh-en-ja-ko-yue-2025-09-09", Label: "Small · Full quality · fp32", Size: "887 MB", Family: sherpa.SenseVoice},
	}
	// moonshineModels are small, fast English recognizers that suit
	// older phones.
	moonshineModels = []model{
		{ID: "sherpa-onnx-moonshine-tiny-en-int8", Label: "Tiny · 27M", Size: "108 MB", Family: sherpa.Moonshine},
		{ID: "sherpa-onnx-moonshine-base-en-int8", Label: "Base · 61M", Size: "251 MB", Family: sherpa.Moonshine},
	}
)

func init() {
	for i := range whisperModels {
		whisperModels[i].Family, whisperModels[i].Int8 = sherpa.Whisper, true
	}
	for _, list := range [][]model{whisperModels, parakeetModels, senseVoiceModels, moonshineModels} {
		for i := range list {
			list[i].Release, list[i].Extra = "asr-models", []string{sileroVAD}
		}
	}
}

// whisperLanguages are the languages Whisper models know; sherpa-onnx
// ends the process when given any other. Cantonese ("yue") came with
// Large v3, so only it and Turbo know it.
var whisperLanguages = strings.Fields(`en zh de es ru ko fr ja pt tr pl ca nl ar sv it id hi fi vi he uk el ms cs ro
	da hu ta no th ur hr bg lt la mi ml cy sk te fa lv bn sr az sl kn et mk br eu is hy ne mn bs kk sq sw gl mr pa si
	km sn yo so af oc ka be tg sd gu am yi lo uz fo ht ps tk nn mt sa lb my bo tl mg as tt haw ln ha ba jw su`)

// senseVoiceLanguages are the language hints SenseVoice accepts.
var senseVoiceLanguages = []string{"zh", "en", "yue", "ja", "ko"}

// supportedLanguage returns lang if model m accepts it as a hint, or ""
// to have the language detected.
func supportedLanguage(m model, lang string) string {
	var known []string
	switch m.Family {
	case sherpa.Whisper:
		known = whisperLanguages
		if m.ID == whisperTurbo || m.ID == whisperLargeV3 {
			known = append(known, "yue")
		}
	case sherpa.SenseVoice:
		known = senseVoiceLanguages
	}
	if slices.Contains(known, lang) {
		return lang
	}
	return ""
}
