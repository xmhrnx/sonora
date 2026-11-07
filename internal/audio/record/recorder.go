package record

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/mehran9675/sonora/internal/audio/opusfile"
)

// Start opens the microphone o.Device and records into o.Path.
func Start(o Options) (*Recorder, error) {
	r, err := newRecorder(o)
	if err != nil {
		return nil, err
	}
	dev, err := open(o.Device, r.push)
	if err != nil {
		r.abort()
		return nil, err
	}
	r.dev = dev
	active.Add(1)
	return r, nil
}

// newRecorder creates the file and starts encoding what push delivers.
func newRecorder(o Options) (*Recorder, error) {
	if err := os.MkdirAll(filepath.Dir(o.Path), 0o755); err != nil {
		return nil, err
	}
	w, err := opusfile.Create(o.Path, SampleRate, o.Bitrate)
	if err != nil {
		return nil, err
	}
	r := &Recorder{path: o.Path, w: w, done: make(chan error, 1)}
	r.wake = sync.NewCond(&r.mu)
	go r.encode()
	return r, nil
}

// Stop ends the recording, finishes the file and returns its path and
// length.
func (r *Recorder) Stop() (string, time.Duration, error) {
	if r.dev != nil {
		r.dev.stop()
		active.Add(-1)
	}
	length := r.Elapsed()
	r.mu.Lock()
	r.closed = true
	r.wake.Signal()
	r.mu.Unlock()
	if err := <-r.done; err != nil {
		os.Remove(r.path)
		return "", 0, err
	}
	return r.path, length, nil
}

// Cancel ends the recording and deletes the file.
func (r *Recorder) Cancel() {
	if path, _, err := r.Stop(); err == nil {
		os.Remove(path)
	}
}

// abort stops encoding and deletes the file, before a device opened.
func (r *Recorder) abort() {
	r.dev = nil
	r.Cancel()
}
