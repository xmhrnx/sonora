//go:build !windows

package platform

// PerMonitorDPI matters only on Windows.
func PerMonitorDPI() {}
