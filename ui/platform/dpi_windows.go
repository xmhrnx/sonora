package platform

import "golang.org/x/sys/windows"

// perMonitorV2 is DPI_AWARENESS_CONTEXT_PER_MONITOR_AWARE_V2.
const perMonitorV2 = ^uintptr(3) // (DPI_AWARENESS_CONTEXT)-4

// PerMonitorDPI keeps the window sharp on monitors with different
// scaling. Sonora's own manifest says so already; a build whose
// manifest only asks for system DPI awareness (the giode CLI's) gets it
// here, before Gio's window asks for less. Call it before the window
// opens. It does nothing on Windows before 10 (1703) or when the
// manifest has set the awareness.
func PerMonitorDPI() {
	proc := windows.NewLazySystemDLL("user32.dll").NewProc("SetProcessDpiAwarenessContext")
	if proc.Find() == nil {
		proc.Call(perMonitorV2)
	}
}
