//go:build sherpa || android

package sherpa

/*
#include <stdlib.h>
#include "c-api.h"
*/
import "C"
import (
	"errors"
	"sync"
	"unsafe"
)

// Recognizer is a loaded speech recognition model. Recognize may be
// called from any goroutine, one call at a time.
type Recognizer struct {
	mu  sync.Mutex
	ptr *C.SherpaOnnxOfflineRecognizer
}

// NewRecognizer loads the model described by cfg.
func NewRecognizer(cfg RecognizerConfig) (*Recognizer, error) {
	var s cstrings
	defer s.free()
	var c C.SherpaOnnxOfflineRecognizerConfig
	c.feat_config.sample_rate = 16000
	c.feat_config.feature_dim = 80
	c.decoding_method = s.str("greedy_search")
	m := &c.model_config
	m.tokens = s.str(cfg.Tokens)
	m.num_threads = threads(cfg.Threads)
	m.provider = s.str("cpu")
	switch cfg.Family {
	case Whisper:
		m.whisper.encoder = s.str(cfg.Encoder)
		m.whisper.decoder = s.str(cfg.Decoder)
		m.whisper.language = s.str(cfg.Language)
		m.whisper.task = s.str("transcribe")
		if cfg.Translate {
			m.whisper.task = s.str("translate")
		}
		m.whisper.tail_paddings = -1
	case Transducer:
		m.transducer.encoder = s.str(cfg.Encoder)
		m.transducer.decoder = s.str(cfg.Decoder)
		m.transducer.joiner = s.str(cfg.Joiner)
		m.model_type = s.str(string(Transducer))
	case SenseVoice:
		m.sense_voice.model = s.str(cfg.Model)
		lang := cfg.Language
		if lang == "" {
			lang = "auto"
		}
		m.sense_voice.language = s.str(lang)
		m.sense_voice.use_itn = 1
	case Moonshine:
		m.moonshine.preprocessor = s.str(cfg.Preprocessor)
		m.moonshine.encoder = s.str(cfg.Encoder)
		m.moonshine.uncached_decoder = s.str(cfg.Decoder)
		m.moonshine.cached_decoder = s.str(cfg.CachedDecoder)
	default:
		return nil, errors.New("unknown model family " + string(cfg.Family))
	}
	ptr := C.SherpaOnnxCreateOfflineRecognizer(&c)
	if ptr == nil {
		return nil, errors.New("the model could not be loaded; its files may be incomplete")
	}
	return &Recognizer{ptr: ptr}, nil
}

// Recognize transcribes one utterance (Whisper takes at most 30 s).
func (r *Recognizer) Recognize(samples []float32, rate int) (Recognition, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ptr == nil {
		return Recognition{}, errors.New("the model was unloaded")
	}
	if len(samples) == 0 {
		return Recognition{}, nil
	}
	stream := C.SherpaOnnxCreateOfflineStream(r.ptr)
	if stream == nil {
		return Recognition{}, errors.New("the recognizer could not start")
	}
	defer C.SherpaOnnxDestroyOfflineStream(stream)
	C.SherpaOnnxAcceptWaveformOffline(stream, C.int32_t(rate), (*C.float)(unsafe.Pointer(&samples[0])), C.int32_t(len(samples)))
	C.SherpaOnnxDecodeOfflineStream(r.ptr, stream)
	res := C.SherpaOnnxGetOfflineStreamResult(stream)
	if res == nil {
		return Recognition{}, errors.New("the recognizer returned nothing")
	}
	defer C.SherpaOnnxDestroyOfflineRecognizerResult(res)
	out := Recognition{Text: C.GoString(res.text)}
	if res.lang != nil {
		out.Language = C.GoString(res.lang)
	}
	return out, nil
}

// Close frees the model. It waits for a running Recognize.
func (r *Recognizer) Close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ptr != nil {
		C.SherpaOnnxDestroyOfflineRecognizer(r.ptr)
		r.ptr = nil
	}
}
