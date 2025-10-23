//go:build linux
// +build linux

package video

import "time"

var (
	cachedPrimaryMonitor *MonitorInfo
	lastPrimaryFetch     time.Time
)

func GetPrimaryMonitorInfo() (*MonitorInfo, error) {
	return nil, ErrOSNotSupported
}

// GetMonitorCount returns the number of connected monitors
func GetMonitorCount() (int, error) {
	return 0, ErrOSNotSupported
}

// GetAllMonitorsInfo returns information about all connected monitors
func GetAllMonitorsInfo() ([]*MonitorInfo, error) {
	return nil, ErrOSNotSupported
}
