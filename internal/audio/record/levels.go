package record

import (
	"math"
	"time"
)

// encode writes queued blocks to the file until the queue is closed.
// The microphone callback only queues, so a slow disk can't make it
// drop audio.
func (r *Recorder) encode() {
	var failed error
	for {
		r.mu.Lock()
		for len(r.queue) == 0 && !r.closed {
			r.wake.Wait()
		}
		blocks, closed := r.queue, r.closed
		r.queue = nil
		r.mu.Unlock()
		for _, b := range blocks {
			if failed == nil {
				failed = r.w.Write(b)
			}
		}
		if closed {
			if err := r.w.Close(); failed == nil {
				failed = err
			}
			r.done <- failed
			return
		}
	}
}

// push receives captured 16-bit little-endian mono PCM from the
// platform device.
func (r *Recorder) push(pcm []byte) {
	n := len(pcm) / 2
	if n == 0 {
		return
	}
	block := make([]float32, n)
	var sum float64
	for i := range block {
		v := float32(int16(uint16(pcm[2*i])|uint16(pcm[2*i+1])<<8)) / 32768
		block[i] = v
		sum += float64(v) * float64(v)
	}
	// Map RMS to a 0..1 meter on a decibel scale (-54 dB .. 0 dB).
	level := float32(0)
	if rms := math.Sqrt(sum / float64(n)); rms > 0 {
		level = float32(max(0, min(1, (20*math.Log10(rms)+54)/54)))
	}
	r.mu.Lock()
	if !r.closed {
		r.queue = append(r.queue, block)
		r.captured += int64(n)
		r.wake.Signal()
	}
	r.levels = append(r.levels, level)
	if len(r.levels) > levelWindow {
		r.levels = r.levels[len(r.levels)-levelWindow:]
	}
	r.mu.Unlock()
}

// Elapsed returns how much audio has been recorded.
func (r *Recorder) Elapsed() time.Duration {
	r.mu.Lock()
	defer r.mu.Unlock()
	return time.Duration(r.captured) * time.Second / SampleRate
}

// Levels returns recent input loudness values in [0, 1], oldest first.
func (r *Recorder) Levels() []float32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return append([]float32(nil), r.levels...)
}
