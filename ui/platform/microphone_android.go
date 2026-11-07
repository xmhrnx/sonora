package platform

import (
	"context"
	"time"

	"gioui.org/app"
	_ "gioui.org/app/permission/microphone" // RECORD_AUDIO in the manifest
	"git.wow.st/gmp/jni"
	"github.com/xmhrnx/giode"
)

// RequestMicrophone asks Android for the microphone permission, showing
// the system prompt the first time, and waits up to a minute for the
// answer. Gio doesn't hand the activity's permission result to apps, so
// the grant is polled. Call it off the UI goroutine.
func RequestMicrophone(a *giode.App, ctx context.Context) error {
	vm := jni.JVMFor(app.JavaVM())
	if micGranted(vm) {
		return nil
	}
	view := a.NativeHandle()
	if view == 0 {
		return errMicDenied
	}
	var err error
	a.Window().Run(func() { err = requestMic(vm, view) })
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()
	t := time.NewTicker(300 * time.Millisecond)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return errMicDenied
		case <-t.C:
			if micGranted(vm) {
				return nil
			}
		}
	}
}
