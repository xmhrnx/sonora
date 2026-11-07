//go:build sherpa || android

package sherpa

/*
#include <stdlib.h>
#include "c-api.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

// VAD is a Silero voice activity detector for 16 kHz audio.
type VAD struct {
	ptr *C.SherpaOnnxVoiceActivityDetector
}

// NewVAD loads silero_vad.onnx from model. Utterances longer than
// maxSpeech seconds are split.
func NewVAD(model string, maxSpeech float32) (*VAD, error) {
	var s cstrings
	defer s.free()
	var c C.SherpaOnnxVadModelConfig
	c.silero_vad.model = s.str(model)
	c.silero_vad.threshold = 0.5
	c.silero_vad.min_silence_duration = 0.4
	c.silero_vad.min_speech_duration = 0.25
	c.silero_vad.window_size = VADWindow
	c.silero_vad.max_speech_duration = C.float(maxSpeech)
	c.sample_rate = 16000
	c.num_threads = 1
	c.provider = s.str("cpu")
	ptr := C.SherpaOnnxCreateVoiceActivityDetector(&c, C.float(maxSpeech+30))
	if ptr == nil {
		return nil, errors.New("the voice activity detector could not be loaded")
	}
	return &VAD{ptr: ptr}, nil
}

// Split feeds 16 kHz samples through the detector and calls fn with
// each stretch of speech and the sample it starts at. fn returning
// false stops early.
func (v *VAD) Split(samples []float32, fn func(start int, speech []float32) bool) {
	C.SherpaOnnxVoiceActivityDetectorReset(v.ptr)
	drain := func() bool {
		for C.SherpaOnnxVoiceActivityDetectorEmpty(v.ptr) == 0 {
			seg := C.SherpaOnnxVoiceActivityDetectorFront(v.ptr)
			speech := make([]float32, int(seg.n))
			copy(speech, unsafe.Slice((*float32)(unsafe.Pointer(seg.samples)), int(seg.n)))
			start := int(seg.start)
			C.SherpaOnnxDestroySpeechSegment(seg)
			C.SherpaOnnxVoiceActivityDetectorPop(v.ptr)
			if !fn(start, speech) {
				return false
			}
		}
		return true
	}
	for i := 0; i+VADWindow <= len(samples); i += VADWindow {
		C.SherpaOnnxVoiceActivityDetectorAcceptWaveform(v.ptr, (*C.float)(unsafe.Pointer(&samples[i])), VADWindow)
		if !drain() {
			return
		}
	}
	C.SherpaOnnxVoiceActivityDetectorFlush(v.ptr)
	drain()
}

// Close frees the detector.
func (v *VAD) Close() {
	if v.ptr != nil {
		C.SherpaOnnxDestroyVoiceActivityDetector(v.ptr)
		v.ptr = nil
	}
}
