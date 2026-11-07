package python

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/mehran9675/sonora/internal/engines"
)

func (w *worker) onStderr(line string) {
	w.mu.Lock()
	w.stderr = append(w.stderr, line)
	if len(w.stderr) > 12 {
		w.stderr = w.stderr[1:]
	}
	log := w.log
	w.mu.Unlock()
	log(line)
}

func (w *worker) setLog(log engines.Logger) {
	w.mu.Lock()
	w.log = log
	w.mu.Unlock()
}

func (w *worker) readLoop(r io.Reader, ready chan<- error) {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 16*1024*1024)
	handshake := false
	for sc.Scan() {
		var m message
		if json.Unmarshal(sc.Bytes(), &m) != nil {
			continue
		}
		switch {
		case m.Event == "ready" && !handshake:
			handshake = true
			ready <- nil
		case m.Event == "fatal" && !handshake:
			handshake = true
			ready <- errors.New(m.Message)
		case m.Event == "log" || m.Event == "progress":
			w.mu.Lock()
			log := w.log
			w.mu.Unlock()
			if m.Event == "progress" {
				log(fmt.Sprintf("%s %d%%", m.Message, int(m.Value*100)))
			} else {
				log(m.Message)
			}
		case m.ID != 0:
			w.mu.Lock()
			ch := w.pending[m.ID]
			delete(w.pending, m.ID)
			w.mu.Unlock()
			if ch != nil {
				ch <- m
			}
		}
	}
}
