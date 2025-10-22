//go:build darwin
// +build darwin

package video

func GetPrimaryMonitorInfo() (*MonitorInfo, error) {
	return nil, ErrOSNotSupported
}
