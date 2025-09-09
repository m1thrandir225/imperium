package input

import (
	"log"
	"syscall"
	"unsafe"
)

var (
	user32        = syscall.MustLoadDLL("user32.dll")
	procSendInput = user32.MustFindProc("SendInput")
)

const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1

	KEYEVENTF_KEYUP    = 0x0002
	KEYEVENTF_SCANCODE = 0x0008

	MOUSEEVENTF_MOVE      = 0x0001
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010

	MOUSEEVENTF_WHEEL = 0x0800
)

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type MOUSEINPUT struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

type INPUT struct {
	Type  uint32
	_     [4]byte // alignment
	Union [24]byte
}

func SendKey(scanCode uint16, keyUp bool) {
	var ki KEYBDINPUT
	ki.WScan = scanCode
	ki.DwFlags = KEYEVENTF_SCANCODE
	if keyUp {
		ki.DwFlags |= KEYEVENTF_KEYUP
	}

	input := INPUT{Type: INPUT_KEYBOARD}
	*(*KEYBDINPUT)(unsafe.Pointer(&input.Union[0])) = ki

	procSendInput.Call(
		1,
		uintptr(unsafe.Pointer(&input)),
		unsafe.Sizeof(input),
	)
}

func SendMouse(flags uint32, dx, dy int32) {
	var mi MOUSEINPUT
	mi.Dx, mi.Dy, mi.DwFlags = dx, dy, flags
	input := INPUT{Type: INPUT_MOUSE}
	*(*MOUSEINPUT)(unsafe.Pointer(&input.Union[0])) = mi

	procSendInput.Call(
		1,
		uintptr(unsafe.Pointer(&input)),
		unsafe.Sizeof(input),
	)
}

func ScrollMouse(amount int32) {
	var mi MOUSEINPUT
	mi.DwFlags = MOUSEEVENTF_WHEEL
	mi.MouseData = uint32(amount)

	input := INPUT{Type: INPUT_MOUSE}
	*(*MOUSEINPUT)(unsafe.Pointer(&input.Union[0])) = mi

	procSendInput.Call(
		1,
		uintptr(unsafe.Pointer(&input)),
		unsafe.Sizeof(input),
	)
}

func KeyToScanCode(key string) uint16 {
	scanCode, ok := keyScancodes[key]
	if !ok {
		return 0
	}
	return scanCode
}

func HandleCommand(cmd InputCommand) {
	switch cmd.Type {
	case "keyboard":
		sc := KeyToScanCode(cmd.Key)
		if sc == 0 {
			log.Printf("Invalid key: %s", cmd.Key)
			return
		}

		if cmd.Action == "press" {
			SendKey(sc, false)
		} else if cmd.Action == "release" {
			SendKey(sc, true)
		}
	case "mouse":
		switch cmd.Action {
		case "move":
			SendMouse(MOUSEEVENTF_MOVE, int32(cmd.X), int32(cmd.Y))
		case "click":
			if cmd.Button == "left" {
				SendMouse(MOUSEEVENTF_LEFTDOWN, 0, 0)
				SendMouse(MOUSEEVENTF_LEFTUP, 0, 0)
			} else if cmd.Button == "right" {
				SendMouse(MOUSEEVENTF_RIGHTDOWN, 0, 0)
				SendMouse(MOUSEEVENTF_RIGHTUP, 0, 0)
			}
		case "scroll":
			ScrollMouse(int32(cmd.Y * 120))
		}
	}
}

var keyScancodes = map[string]uint16{
	"A":            0x1E,
	"B":            0x30,
	"C":            0x2E,
	"D":            0x20,
	"E":            0x12,
	"F":            0x21,
	"G":            0x22,
	"H":            0x23,
	"I":            0x17,
	"J":            0x24,
	"K":            0x25,
	"L":            0x26,
	"M":            0x32,
	"N":            0x31,
	"O":            0x18,
	"P":            0x19,
	"Q":            0x10,
	"R":            0x13,
	"S":            0x1F,
	"T":            0x14,
	"U":            0x16,
	"V":            0x2F,
	"W":            0x11,
	"X":            0x2D,
	"Y":            0x15,
	"Z":            0x2C,
	"1":            0x02,
	"2":            0x03,
	"3":            0x04,
	"4":            0x05,
	"5":            0x06,
	"6":            0x07,
	"7":            0x08,
	"8":            0x09,
	"9":            0x0A,
	"0":            0x0B,
	"-":            0x0C,
	"=":            0x0D,
	"`":            0x0E,
	"[":            0x1A,
	"]":            0x1B,
	"\\":           0x2B,
	"#":            0x27,
	"*":            0x37,
	"+":            0x0F,
	"/":            0x39,
	"CapsLock":     0x3A,
	"F1":           0x3B,
	"F2":           0x3C,
	"F3":           0x3D,
	"F4":           0x3E,
	"F5":           0x3F,
	"F6":           0x40,
	"F7":           0x41,
	"F8":           0x42,
	"F9":           0x43,
	"F10":          0x44,
	"F11":          0x57,
	"F12":          0x58,
	"Insert":       0x49,
	"Delete":       0x53,
	"Home":         0x47,
	"End":          0x4F,
	"PageUp":       0x4B,
	"PageDown":     0x4E,
	"Up":           0xC8,
	"Down":         0xD0,
	"Left":         0xCB,
	"Right":        0xD3,
	"Tab":          0x0F,
	"Enter":        0x1C,
	"Escape":       0x01,
	"Backspace":    0x0E,
	"Space":        0x39,
	"Shift":        0x2A,
	"Ctrl":         0x1D,
	"Alt":          0x38,
	"Win":          0x5B,
	"Apps":         0x5D,
	"NumLock":      0x45,
	"ScrollLock":   0x46,
	"Pause":        0x48,
	"PrintScreen":  0x46,
	"Sleep":        0x48,
	"Num0":         0x52,
	"Num1":         0x4F,
	"Num2":         0x50,
	"Num3":         0x51,
	"Num4":         0x4B,
	"Num5":         0x4C,
	"Num6":         0x4D,
	"Num7":         0x47,
	"Num8":         0x48,
	"Num9":         0x49,
	"NumAdd":       0x4E,
	"NumSubtract":  0x4A,
	"NumMultiply":  0x37,
	"NumDivide":    0x38,
	"NumDecimal":   0x53,
	"NumEnter":     0x1C,
	"NumBackspace": 0x0E,
	"NumTab":       0x0F,
	"NumHome":      0x47,
	"NumEnd":       0x4F,
	"NumPageUp":    0x4B,
	"NumPageDown":  0x4E,
}
