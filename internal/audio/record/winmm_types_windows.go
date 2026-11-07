package record

import (
	"sync"

	"golang.org/x/sys/windows"
)

type waveFormatEx struct {
	FormatTag      uint16
	Channels       uint16
	SamplesPerSec  uint32
	AvgBytesPerSec uint32
	BlockAlign     uint16
	BitsPerSample  uint16
	Size           uint16
}

// waveInCaps is WAVEINCAPSW.
type waveInCaps struct {
	Mid, Pid      uint16
	DriverVersion uint32
	Pname         [32]uint16
	Formats       uint32
	Channels      uint16
	Reserved      uint16
}

type waveHdr struct {
	Data          uintptr
	BufferLength  uint32
	BytesRecorded uint32
	User          uintptr
	Flags         uint32
	Loops         uint32
	Next          uintptr
	Reserved      uintptr
}

type waveIn struct {
	handle  uintptr
	event   windows.Handle
	hdrs    [bufferCount]*waveHdr
	bufs    [bufferCount][]byte
	push    func([]byte)
	quit    chan struct{}
	done    sync.WaitGroup
	stopped bool
}

// winmm captures with the classic waveIn API, which every Windows
// version has and which needs no cgo, so ARM64 builds stay simple.
var (
	winmm                = windows.NewLazySystemDLL("winmm.dll")
	procWaveInGetNumDevs = winmm.NewProc("waveInGetNumDevs")
	procWaveInGetDevCaps = winmm.NewProc("waveInGetDevCapsW")
	procWaveInOpen       = winmm.NewProc("waveInOpen")
	procWaveInPrepare    = winmm.NewProc("waveInPrepareHeader")
	procWaveInUnprepare  = winmm.NewProc("waveInUnprepareHeader")
	procWaveInAddBuffer  = winmm.NewProc("waveInAddBuffer")
	procWaveInStart      = winmm.NewProc("waveInStart")
	procWaveInStop       = winmm.NewProc("waveInStop")
	procWaveInReset      = winmm.NewProc("waveInReset")
	procWaveInClose      = winmm.NewProc("waveInClose")
)

const (
	waveMapper    = 0xFFFFFFFF
	callbackEvent = 0x00050000
	whdrDone      = 0x00000001
	bufferCount   = 4
	bufferBytes   = SampleRate / 10 * 2 // 100 ms of 16-bit mono
)
