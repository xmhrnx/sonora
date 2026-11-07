package onnx

import "github.com/mehran9675/sonora/internal/engines"

func whisperManifest() *engines.Manifest {
	return &engines.Manifest{
		ID: "ondevice-whisper", Name: "Whisper", Vendor: "OpenAI · on device",
		Description: "Transcribes speech in about 100 languages, offline on this device, and can translate it into English. " +
			"Bigger models are more accurate and slower; Turbo suits flagship phones, Medium and Large v3 computers.",
		Tags:    []string{"Offline", "99 languages", "Translate"},
		License: "MIT", Homepage: "https://github.com/openai/whisper",
		Size: "110 MB–1.9 GB", Accent: []string{"#10b981", "#0d9488"},
		Kind: engines.KindOnDevice, Task: engines.TaskTranscribe, Translate: true,
		Options: []engines.Option{modelOption(whisperModels, "A new model downloads the first time you use it.")},
	}
}

func parakeetManifest() *engines.Manifest {
	return &engines.Manifest{
		ID: "ondevice-parakeet", Name: "Parakeet", Vendor: "NVIDIA · on device",
		Description: "NVIDIA's fast, accurate recognizer for English and 24 other European languages, with punctuation, " +
			"offline on this device. Needs about 1 GB of free memory.",
		Tags:    []string{"Offline", "25 languages", "Accurate"},
		License: "CC-BY-4.0", Homepage: "https://huggingface.co/nvidia/parakeet-tdt-0.6b-v3",
		Size: "464–482 MB", Accent: []string{"#76b900", "#16a34a"},
		// Parakeet detects which of its languages is spoken and can't be
		// told.
		Kind: engines.KindOnDevice, Task: engines.TaskTranscribe, Languages: []string{""},
		Options: []engines.Option{modelOption(parakeetModels, "v3 understands 25 European languages; v2 only English.")},
	}
}

func senseVoiceManifest() *engines.Manifest {
	return &engines.Manifest{
		ID: "ondevice-sensevoice", Name: "SenseVoice", Vendor: "FunAudioLLM · on device",
		Description: "Very fast recognition of Chinese, Cantonese, English, Japanese and Korean, offline on this device.",
		Tags:        []string{"Offline", "Fast", "Chinese · English · Japanese · Korean"},
		License:     "FunASR model license", Homepage: "https://github.com/FunAudioLLM/SenseVoice",
		Size: "158–887 MB", Accent: []string{"#f59e0b", "#ef4444"},
		Kind: engines.KindOnDevice, Task: engines.TaskTranscribe, Languages: append([]string{""}, senseVoiceLanguages...),
		Options: []engines.Option{modelOption(senseVoiceModels, "Full quality is a little more accurate and five times the download.")},
	}
}

func moonshineManifest() *engines.Manifest {
	return &engines.Manifest{
		ID: "ondevice-moonshine", Name: "Moonshine", Vendor: "Useful Sensors · on device",
		Description: "Small, fast English speech recognition that runs offline even on older phones.",
		Tags:        []string{"Offline", "Fast", "English"},
		License:     "MIT", Homepage: "https://github.com/usefulsensors/moonshine",
		Size: "108–251 MB", Accent: []string{"#fbbf24", "#a855f7"},
		Kind: engines.KindOnDevice, Task: engines.TaskTranscribe, Languages: []string{"", "en"},
		Options: []engines.Option{modelOption(moonshineModels, "Base is more accurate; Tiny is faster.")},
	}
}

func barkManifest() *engines.Manifest {
	return &engines.Manifest{
		ID: "ondevice-bark", Name: "Underdog Bark", Vendor: "Conway Research · on device",
		Description: "Conway Research's English speech recognition, tuned from Parakeet to be more accurate, offline on this device. " +
			"It writes lowercase text without punctuation. Plus is a little more accurate.",
		Tags:     []string{"Offline", "English", "Accurate"},
		License:  "CC-BY-4.0 (converted weights) · Conway's own release is CC-BY-NC-4.0",
		Homepage: "https://huggingface.co/ConwayResearch/Underdog-Bark-3.0",
		Size:     "397–685 MB", Accent: []string{"#f97316", "#a16207"},
		Kind: engines.KindOnDevice, Task: engines.TaskTranscribe, Languages: []string{"", "en"},
		Options: []engines.Option{modelOption(barkModels, "Compact downloads less; Plus is the most accurate.")},
	}
}
