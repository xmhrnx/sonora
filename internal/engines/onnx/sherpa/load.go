//go:build sherpa || android

package sherpa

/*
#cgo CFLAGS: -I${SRCDIR}
#cgo LDFLAGS: -lsherpa-onnx-c-api -lonnxruntime
#cgo android,arm64 LDFLAGS: -L${SRCDIR}/../../../../third_party/sherpa-onnx/android/arm64-v8a
#cgo android,arm LDFLAGS: -L${SRCDIR}/../../../../third_party/sherpa-onnx/android/armeabi-v7a
#cgo android,amd64 LDFLAGS: -L${SRCDIR}/../../../../third_party/sherpa-onnx/android/x86_64
#cgo linux,!android,amd64 LDFLAGS: -L${SRCDIR}/../../../../third_party/sherpa-onnx/linux-x64
#include <stdlib.h>
#include <string.h>
#include "c-api.h"

extern int32_t sonoraProgress(float *samples, int32_t n, float p, void *arg);

static const SherpaOnnxGeneratedAudio *generate(const SherpaOnnxOfflineTts *tts,
		const char *text, SherpaOnnxGenerationConfig *cfg, uintptr_t handle) {
	cfg->silence_scale = 0.2f;
	return SherpaOnnxOfflineTtsGenerateWithConfig(tts, text, cfg,
		(SherpaOnnxGeneratedAudioProgressCallbackWithArg)sonoraProgress, (void *)handle);
}
*/
import "C"
import (
	"errors"
	"unsafe"
)

// New loads the model described by cfg.
func New(cfg Config) (*TTS, error) {
	var strs []*C.char
	str := func(s string) *C.char {
		if s == "" {
			return nil
		}
		c := C.CString(s)
		strs = append(strs, c)
		return c
	}
	defer func() {
		for _, c := range strs {
			C.free(unsafe.Pointer(c))
		}
	}()
	var c C.SherpaOnnxOfflineTtsConfig
	c.model.num_threads = threads(cfg.Threads)
	c.model.provider = str("cpu")
	c.max_num_sentences = 1
	switch cfg.Family {
	case VITS:
		c.model.vits.model = str(cfg.Model)
		c.model.vits.tokens = str(cfg.Tokens)
		c.model.vits.data_dir = str(cfg.DataDir)
		c.model.vits.noise_scale = 0.667
		c.model.vits.noise_scale_w = 0.8
		c.model.vits.length_scale = 1
	case Kitten:
		c.model.kitten.model = str(cfg.Model)
		c.model.kitten.voices = str(cfg.Voices)
		c.model.kitten.tokens = str(cfg.Tokens)
		c.model.kitten.data_dir = str(cfg.DataDir)
		c.model.kitten.length_scale = 1
	case Kokoro:
		c.model.kokoro.model = str(cfg.Model)
		c.model.kokoro.voices = str(cfg.Voices)
		c.model.kokoro.tokens = str(cfg.Tokens)
		c.model.kokoro.data_dir = str(cfg.DataDir)
		c.model.kokoro.lexicon = str(cfg.Lexicon)
		c.model.kokoro.dict_dir = str(cfg.DictDir)
		c.model.kokoro.length_scale = 1
	case ZipVoice:
		c.model.zipvoice.tokens = str(cfg.Tokens)
		c.model.zipvoice.encoder = str(cfg.Encoder)
		c.model.zipvoice.decoder = str(cfg.Decoder)
		c.model.zipvoice.vocoder = str(cfg.Vocoder)
		c.model.zipvoice.data_dir = str(cfg.DataDir)
		c.model.zipvoice.lexicon = str(cfg.Lexicon)
		c.model.zipvoice.feat_scale = 0.1
		c.model.zipvoice.t_shift = 0.5
		c.model.zipvoice.target_rms = 0.1
		c.model.zipvoice.guidance_scale = 1
	default:
		return nil, errors.New("unknown model family " + string(cfg.Family))
	}
	ptr := C.SherpaOnnxCreateOfflineTts(&c)
	if ptr == nil {
		return nil, errors.New("the model could not be loaded; its files may be incomplete")
	}
	return &TTS{ptr: ptr}, nil
}

// SampleRate returns the model's output rate.
func (t *TTS) SampleRate() int { return int(C.SherpaOnnxOfflineTtsSampleRate(t.ptr)) }

// NumSpeakers returns how many voices the model has.
func (t *TTS) NumSpeakers() int { return int(C.SherpaOnnxOfflineTtsNumSpeakers(t.ptr)) }
