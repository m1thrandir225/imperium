//go:build windows
// +build windows

package programs

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	user32              = syscall.NewLazyDLL("user32.dll")
	_                   = syscall.NewLazyDLL("kernel32.dll")
	procEnumWindows     = user32.NewProc("EnumWindows")
	procGetWindowThread = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowTextW  = user32.NewProc("GetWindowTextW")
	procIsWindowVisible = user32.NewProc("IsWindowVisible")
)

func (s *programService) GetWindowTitleByProcessID(pid uint32) (string, error) {
	var hwnd syscall.Handle
	cb := syscall.NewCallback(func(h syscall.Handle, lparam uintptr) uintptr {
		var processID uint32
		procGetWindowThread.Call(uintptr(h), uintptr(unsafe.Pointer(&processID)))
		if processID == pid {
			// Check if window is visible
			ret, _, _ := procIsWindowVisible.Call(uintptr(h))
			if ret == 0 {
				return 1 // continue
			}

			// Get window text
			buf := make([]uint16, 256)
			procGetWindowTextW.Call(uintptr(h), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
			title := syscall.UTF16ToString(buf)
			if title != "" {
				hwnd = h
				return 0 // stop
			}
		}
		return 1 // continue
	})

	procEnumWindows.Call(cb, 0)

	if hwnd == 0 {
		return "", fmt.Errorf("no window found for process ID %d", pid)
	}

	// Extract window text again
	buf := make([]uint16, 256)
	procGetWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	return syscall.UTF16ToString(buf), nil
}
