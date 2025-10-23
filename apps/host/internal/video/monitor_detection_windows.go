//go:build windows
// +build windows

package video

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

// Windows API constants
const (
	MONITOR_DEFAULTTOPRIMARY = 1
	MONITOR_DEFAULTTONEAREST = 2
	MONITOR_DEFAULTTONULL    = 0
)

// Windows API structures
type RECT struct {
	Left, Top, Right, Bottom int32
}

type MONITORINFO struct {
	CbSize    uint32
	RcMonitor RECT
	RcWork    RECT
	DwFlags   uint32
}

// Windows API functions
var (
	user32                  = syscall.NewLazyDLL("user32.dll")
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetMonitorInfoW     = user32.NewProc("GetMonitorInfoW")
	procMonitorFromPoint    = user32.NewProc("MonitorFromPoint")
)

var (
	cachedPrimaryMonitor *MonitorInfo
	lastPrimaryFetch     time.Time
)

// GetPrimaryMonitorInfo returns information about the primary monitor
func GetPrimaryMonitorInfo() (*MonitorInfo, error) {
	if cachedPrimaryMonitor != nil && time.Since(lastPrimaryFetch) < 10*time.Minute {
		return cachedPrimaryMonitor, nil
	}

	var primaryMonitor *MonitorInfo

	// Callback function for EnumDisplayMonitors
	enumProc := syscall.NewCallback(func(hMonitor syscall.Handle, hdc syscall.Handle, lprcClip *RECT, dwData uintptr) uintptr {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		ret, _, _ := procGetMonitorInfoW.Call(
			uintptr(hMonitor),
			uintptr(unsafe.Pointer(&mi)),
		)

		if ret != 0 {
			// Check if this is the primary monitor
			isPrimary := (mi.DwFlags & 1) != 0 // MONITORINFOF_PRIMARY = 1

			if isPrimary {
				primaryMonitor = &MonitorInfo{
					Width:     int(mi.RcMonitor.Right - mi.RcMonitor.Left),
					Height:    int(mi.RcMonitor.Bottom - mi.RcMonitor.Top),
					OffsetX:   int(mi.RcMonitor.Left),
					OffsetY:   int(mi.RcMonitor.Top),
					IsPrimary: true,
				}
			}
		}

		return 1 // Continue enumeration
	})

	// Enumerate all monitors
	ret, _, _ := procEnumDisplayMonitors.Call(
		0, // hdc
		0, // lprcClip
		enumProc,
		0, // dwData
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate monitors")
	}

	if primaryMonitor == nil {
		return nil, fmt.Errorf("primary monitor not found")
	}
	cachedPrimaryMonitor = primaryMonitor
	lastPrimaryFetch = time.Now()

	return primaryMonitor, nil
}

// GetMonitorCount returns the number of connected monitors
func GetMonitorCount() (int, error) {
	count := 0

	// Callback function for EnumDisplayMonitors
	enumProc := syscall.NewCallback(func(hMonitor syscall.Handle, hdc syscall.Handle, lprcClip *RECT, dwData uintptr) uintptr {
		count++
		return 1 // Continue enumeration
	})

	// Enumerate all monitors
	ret, _, _ := procEnumDisplayMonitors.Call(
		0, // hdc
		0, // lprcClip
		enumProc,
		0, // dwData
	)

	if ret == 0 {
		return 0, fmt.Errorf("failed to enumerate monitors")
	}

	return count, nil
}

// GetAllMonitorsInfo returns information about all connected monitors on Windows
func GetAllMonitorsInfo() ([]*MonitorInfo, error) {
	var monitors []*MonitorInfo

	// Callback function for EnumDisplayMonitors
	enumProc := syscall.NewCallback(func(hMonitor syscall.Handle, hdc syscall.Handle, lprcClip *RECT, dwData uintptr) uintptr {
		var mi MONITORINFO
		mi.CbSize = uint32(unsafe.Sizeof(mi))

		ret, _, _ := procGetMonitorInfoW.Call(
			uintptr(hMonitor),
			uintptr(unsafe.Pointer(&mi)),
		)

		if ret != 0 {
			// Check if this is the primary monitor
			isPrimary := (mi.DwFlags & 1) != 0 // MONITORINFOF_PRIMARY = 1

			monitor := &MonitorInfo{
				Width:     int(mi.RcMonitor.Right - mi.RcMonitor.Left),
				Height:    int(mi.RcMonitor.Bottom - mi.RcMonitor.Top),
				OffsetX:   int(mi.RcMonitor.Left),
				OffsetY:   int(mi.RcMonitor.Top),
				IsPrimary: isPrimary,
			}

			monitors = append(monitors, monitor)
		}

		return 1 // Continue enumeration
	})

	ret, _, _ := procEnumDisplayMonitors.Call(
		0, // hdc
		0, // lprcClip
		enumProc,
		0, // dwData
	)

	if ret == 0 {
		return nil, fmt.Errorf("failed to enumerate monitors")
	}

	return monitors, nil
}
