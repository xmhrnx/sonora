package python

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// call sends a command and waits for its reply. Cancelling ctx kills
// the worker, since a model mid-generation can't be interrupted.
func (w *worker) call(ctx context.Context, cmd string, opts map[string]any, req any) (json.RawMessage, error) {
	w.mu.Lock()
	w.nextID++
	id := w.nextID
	ch := make(chan message, 1)
	w.pending[id] = ch
	w.mu.Unlock()
	line, _ := json.Marshal(request{ID: id, Cmd: cmd, Options: opts, Request: req})
	if _, err := w.stdin.Write(append(line, '\n')); err != nil {
		return nil, w.failure("the engine is not running")
	}
	select {
	case m := <-ch:
		if !m.OK {
			return nil, errors.New(m.Error)
		}
		return m.Result, nil
	case <-w.dead:
		return nil, w.failure("the engine stopped unexpectedly")
	case <-ctx.Done():
		w.close()
		return nil, ctx.Err()
	}
}

func (w *worker) alive() bool {
	select {
	case <-w.dead:
		return false
	default:
		return true
	}
}

// failure builds an error from the last stderr lines, which usually
// hold the Python traceback.
func (w *worker) failure(what string) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if len(w.stderr) == 0 {
		return errors.New(what)
	}
	return fmt.Errorf("%s:\n%s", what, strings.Join(w.stderr, "\n"))
}

func (w *worker) close() {
	if !w.alive() {
		return
	}
	w.stdin.Write([]byte(`{"cmd":"quit"}` + "\n"))
	w.stdin.Close()
	select {
	case <-w.dead:
	case <-time.After(2 * time.Second):
		w.cmd.Process.Kill()
	}
}
