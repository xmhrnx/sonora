package onnx

import "github.com/mehran9675/sonora/internal/engines/onnx/sherpa"

var (
	kittenModels = []model{
		{ID: "kitten-nano-en-v0_8-int8", Label: "Nano · 15M", Size: "31 MB", Family: sherpa.Kitten},
		{ID: "kitten-micro-en-v0_8", Label: "Micro · 40M", Size: "44 MB", Family: sherpa.Kitten},
		{ID: "kitten-mini-en-v0_8", Label: "Mini · 80M", Size: "68 MB", Family: sherpa.Kitten},
	}
	// kokoroModels share one voice list; the full-precision model
	// sounds a little cleaner and needs about 1 GB of free memory.
	kokoroModels = []model{
		{ID: "kokoro-int8-multi-lang-v1_0", Label: "Compact · int8", Size: "132 MB",
			Family: sherpa.Kokoro, Lexicon: "lexicon-us-en.txt", Dict: "dict"},
		{ID: "kokoro-multi-lang-v1_0", Label: "Full quality · fp32", Size: "333 MB",
			Family: sherpa.Kokoro, Lexicon: "lexicon-us-en.txt", Dict: "dict"},
	}
	// zipvoiceModels clone a voice from a recording and its transcript.
	// Both are the distilled model, which needs only 4 steps; sizes
	// include the shared 52 MB vocoder.
	zipvoiceModels = []model{
		{ID: "sherpa-onnx-zipvoice-distill-int8-zh-en-emilia", Label: "Compact · int8", Size: "156 MB",
			Family: sherpa.ZipVoice, Lexicon: "lexicon.txt", Steps: 4, Extra: []string{vocos}},
		{ID: "sherpa-onnx-zipvoice-distill-fp32-zh-en-emilia", Label: "Full quality · fp32", Size: "508 MB",
			Family: sherpa.ZipVoice, Lexicon: "lexicon.txt", Steps: 4, Extra: []string{vocos}},
	}
)

// kittenVoices are speaker ids 0-7 in sherpa-onnx's voices.bin
// (expr-voice-2-m, -2-f, -3-m, -3-f, …), with KittenML's names.
var kittenVoices = []struct{ Name, Gender string }{
	{"Jasper", "Male"}, {"Bella", "Female"}, {"Bruno", "Male"}, {"Luna", "Female"},
	{"Hugo", "Male"}, {"Rosie", "Female"}, {"Leo", "Male"}, {"Kiki", "Female"},
}

// kokoroVoices are Kokoro v1.0's English speakers; their position is
// the speaker id in sherpa-onnx's voices.bin.
var kokoroVoices = []string{
	"af_alloy", "af_aoede", "af_bella", "af_heart", "af_jessica", "af_kore", "af_nicole", "af_nova",
	"af_river", "af_sarah", "af_sky", "am_adam", "am_echo", "am_eric", "am_fenrir", "am_liam",
	"am_michael", "am_onyx", "am_puck", "am_santa", "bf_alice", "bf_emma", "bf_isabella", "bf_lily",
	"bm_daniel", "bm_fable", "bm_george", "bm_lewis",
}
