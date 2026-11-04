//go:build !android

// Package platform holds what differs between operating systems: the
// microphone and its permission.
package platform

import (
	"context"

	"github.com/xmhrnx/giode"

	"github.com/mehran9675/sonora/internal/audio/record"
)

// RequestMicrophone asks for access to the microphone where the
// platform needs an explicit grant (Android); elsewhere it's a no-op.
func RequestMicrophone(*giode.App, context.Context) error { return nil }

// Microphones lists the microphones to record from. It can take a
// moment; call it off the UI goroutine.
func Microphones() ([]record.Device, error) { return record.Devices() }
