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
	"runtime/cgo"
	"unsafe"
)

type genState struct {
	progress Progress
	stopped  bool
}

// Generate speaks text with the settings in p, reporting progress.
func (t *TTS) Generate(text string, p Params, progress Progress) ([]float32, int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ptr == nil {
		return nil, 0, errors.New("the model was unloaded")
	}
	ct := C.CString(text)
	defer C.free(unsafe.Pointer(ct))
	var cfg C.SherpaOnnxGenerationConfig
	cfg.sid = C.int32_t(p.Speaker)
	cfg.speed = C.float(p.Speed)
	cfg.num_steps = C.int32_t(p.Steps)
	if len(p.Reference) > 0 {
		// Copied to C memory: cgo forbids passing Go pointers inside a
		// struct.
		ref := (*C.float)(C.malloc(C.size_t(len(p.Reference)) * C.size_t(unsafe.Sizeof(C.float(0)))))
		defer C.free(unsafe.Pointer(ref))
		copy(unsafe.Slice((*float32)(unsafe.Pointer(ref)), len(p.Reference)), p.Reference)
		cfg.reference_audio = ref
		cfg.reference_audio_len = C.int32_t(len(p.Reference))
		cfg.reference_sample_rate = C.int32_t(p.ReferenceRate)
	}
	if p.ReferenceText != "" {
		rt := C.CString(p.ReferenceText)
		defer C.free(unsafe.Pointer(rt))
		cfg.reference_text = rt
	}
	state := &genState{progress: progress}
	h := cgo.NewHandle(state)
	defer h.Delete()
	audio := C.generate(t.ptr, ct, &cfg, C.uintptr_t(h))
	if state.stopped {
		if audio != nil {
			C.SherpaOnnxDestroyOfflineTtsGeneratedAudio(audio)
		}
		return nil, 0, ErrStopped
	}
	if audio == nil {
		return nil, 0, errors.New("the model produced no audio")
	}
	defer C.SherpaOnnxDestroyOfflineTtsGeneratedAudio(audio)
	n := int(audio.n)
	out := make([]float32, n)
	if n > 0 {
		copy(out, unsafe.Slice((*float32)(unsafe.Pointer(audio.samples)), n))
	}
	return out, int(audio.sample_rate), nil
}

// Close frees the model. It waits for a running Generate.
func (t *TTS) Close() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ptr != nil {
		C.SherpaOnnxDestroyOfflineTts(t.ptr)
		t.ptr = nil
	}
}

//export sonoraProgress
func sonoraProgress(samples *C.float, n C.int32_t, p C.float, arg unsafe.Pointer) C.int32_t {
	s := cgo.Handle(uintptr(arg)).Value().(*genState)
	if s.progress != nil && !s.progress(float32(p)) {
		s.stopped = true
		return 0
	}
	return 1
}
