//go:build linux
// +build linux

package video

func GetPrimaryMonitorInfo() (*MonitorInfo, error) {
	return nil, ErrOSNotSupported
}
