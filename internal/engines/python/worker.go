package python

import (
	"context"
	"io"
	"os/exec"
	"sync"

	"github.com/mehran9675/sonora/internal/engines"
	"github.com/mehran9675/sonora/internal/paths"
	"github.com/mehran9675/sonora/internal/pyenv"
)

// worker is a running worker.py process hosting one adapter.
type worker struct {
	cmd   *exec.Cmd
	stdin io.WriteCloser

	mu      sync.Mutex
	nextID  int
	pending map[int]chan message
	log     engines.Logger
	dead    chan struct{}
	exitErr error
	stderr  []string
}

// startWorker launches the worker and waits for its handshake.
func startWorker(ctx context.Context, python, script, adapter, engineDir string, env map[string]string, log engines.Logger) (*worker, error) {
	cmd := exec.Command(python, "-u", script, "--adapter", adapter, "--engine-dir", engineDir)
	cmd.Env = pyenv.Env(env)
	paths.HideWindow(cmd)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	w := &worker{cmd: cmd, stdin: stdin, pending: map[int]chan message{}, log: log, dead: make(chan struct{})}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	ready := make(chan error, 1)
	go pyenv.ScanLines(stderr, w.onStderr)
	go w.readLoop(stdout, ready)
	go func() {
		w.exitErr = cmd.Wait()
		close(w.dead)
	}()
	select {
	case err := <-ready:
		if err != nil {
			w.close()
			return nil, err
		}
		return w, nil
	case <-w.dead:
		return nil, w.failure("the engine stopped while starting")
	case <-ctx.Done():
		w.close()
		return nil, ctx.Err()
	}
}
