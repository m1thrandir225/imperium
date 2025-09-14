package input

import (
	"fmt"
	"log"
	"strings"

	"github.com/stephen-fox/user32util"
)

var user32DLL *user32util.User32DLL

func init() {
	var err error
	user32DLL, err = user32util.LoadUser32DLL()
	if err != nil {
		log.Fatalf("Failed to load user32 DLL: %v", err)
	}
	log.Printf("✅ Successfully loaded user32 DLL")
}

func HandleCommand(cmd InputCommand) {
	log.Printf("🎮 Processing input command: Type=%s, Action=%s, Key=%s, X=%d, Y=%d, Button=%s",
		cmd.Type, cmd.Action, cmd.Key, cmd.X, cmd.Y, cmd.Button)

	switch cmd.Type {
	case "mouse":
		handleMouseCommand(cmd)
	case "keyboard":
		handleKeyboardCommand(cmd)
	default:
		log.Printf("❌ Unknown input type: %s", cmd.Type)
	}
}

func handleMouseCommand(cmd InputCommand) {
	switch cmd.Action {
	case "move":
		log.Printf("🖱️ Moving mouse: X=%d, Y=%d", cmd.X, cmd.Y)
		if err := moveMouse(cmd.X, cmd.Y); err != nil {
			log.Printf("❌ Failed to move mouse: %v", err)
		} else {
			log.Printf("✅ Mouse moved successfully")
		}
	case "press":
		log.Printf("🖱️ Mouse press: %s", cmd.Button)
		if err := pressMouseButton(cmd.Button); err != nil {
			log.Printf("❌ Failed to press mouse button: %v", err)
		} else {
			log.Printf("✅ Mouse button pressed successfully")
		}
	case "release":
		log.Printf("🖱️ Mouse release: %s", cmd.Button)
		if err := releaseMouseButton(cmd.Button); err != nil {
			log.Printf("❌ Failed to release mouse button: %v", err)
		} else {
			log.Printf("✅ Mouse button released successfully")
		}
	case "click":
		log.Printf("🖱️ Mouse click: %s", cmd.Button)
		if err := clickMouseButton(cmd.Button); err != nil {
			log.Printf("❌ Failed to click mouse button: %v", err)
		} else {
			log.Printf("✅ Mouse clicked successfully")
		}
	default:
		log.Printf("❌ Unknown mouse action: %s", cmd.Action)
	}
}

func handleKeyboardCommand(cmd InputCommand) {
	switch cmd.Action {
	case "press":
		log.Printf("⌨️ Key press: %s", cmd.Key)
		if err := pressKey(cmd.Key); err != nil {
			log.Printf("❌ Failed to press key: %v", err)
		} else {
			log.Printf("✅ Key pressed successfully")
		}
	case "release":
		log.Printf("⌨️ Key release: %s", cmd.Key)
		if err := releaseKey(cmd.Key); err != nil {
			log.Printf("❌ Failed to release key: %v", err)
		} else {
			log.Printf("✅ Key released successfully")
		}
	default:
		log.Printf("❌ Unknown keyboard action: %s", cmd.Action)
	}
}

func moveMouse(x, y int) error {
	log.Printf("🔍 Setting cursor position: X=%d, Y=%d", x, y)

	_, err := user32util.SetCursorPos(int32(x), int32(y), user32DLL)
	if err != nil {
		return fmt.Errorf("SetCursorPos failed: %w", err)
	}

	return nil
}

func pressMouseButton(button string) error {
	var flags uint32
	switch strings.ToLower(button) {
	case "left":
		flags = user32util.MouseEventFLeftDown
	case "right":
		flags = user32util.MouseEventFRightDown
	case "middle":
		flags = user32util.MouseEventFMiddleDown
	default:
		return fmt.Errorf("unknown mouse button: %s", button)
	}

	mouseInput := user32util.MouseInput{
		DwFlags: flags,
	}

	log.Printf("🔍 Sending mouse button press: %s (flags: %v)", button, flags)

	err := user32util.SendMouseInput(mouseInput, user32DLL)
	if err != nil {
		return fmt.Errorf("SendMouseInput failed: %w", err)
	}

	return nil
}

func releaseMouseButton(button string) error {
	var flags uint32
	switch strings.ToLower(button) {
	case "left":
		flags = user32util.MouseEventFLeftUp
	case "right":
		flags = user32util.MouseEventFRightUp
	case "middle":
		flags = user32util.MouseEventFMiddleUp
	default:
		return fmt.Errorf("unknown mouse button: %s", button)
	}

	mouseInput := user32util.MouseInput{
		DwFlags: flags,
	}

	log.Printf("🔍 Sending mouse button release: %s (flags: %v)", button, flags)

	err := user32util.SendMouseInput(mouseInput, user32DLL)
	if err != nil {
		return fmt.Errorf("SendMouseInput failed: %w", err)
	}

	return nil
}

func clickMouseButton(button string) error {
	if err := pressMouseButton(button); err != nil {
		return err
	}
	if err := releaseMouseButton(button); err != nil {
		return err
	}
	return nil
}

var keyMap = map[string]uint16{
	"a": 0x41, "b": 0x42, "c": 0x43, "d": 0x44, "e": 0x45, "f": 0x46, "g": 0x47, "h": 0x48,
	"i": 0x49, "j": 0x4A, "k": 0x4B, "l": 0x4C, "m": 0x4D, "n": 0x4E, "o": 0x4F, "p": 0x50,
	"q": 0x51, "r": 0x52, "s": 0x53, "t": 0x54, "u": 0x55, "v": 0x56, "w": 0x57, "x": 0x58,
	"y": 0x59, "z": 0x5A,
	"0": 0x30, "1": 0x31, "2": 0x32, "3": 0x33, "4": 0x34, "5": 0x35, "6": 0x36, "7": 0x37, "8": 0x38, "9": 0x39,
	" ":     0x20, // Space character
	"":      0x20, // Empty string (fallback for space)
	"space": 0x20, // Keep the existing space mapping too
	"enter": 0x0D, "escape": 0x1B, "tab": 0x09, "shift": 0x10, "ctrl": 0x11, "alt": 0x12,
	"left": 0x25, "up": 0x26, "right": 0x27, "down": 0x28,
	"arrowleft": 0x25, "arrowup": 0x26, "arrowright": 0x27, "arrowdown": 0x28,
	"meta": 0x20,
}

func pressKey(key string) error {
	vkCode, exists := keyMap[strings.ToLower(key)]
	if !exists {
		return fmt.Errorf("unknown key: %s", key)
	}

	keyInput := user32util.KeybdInput{
		WVK: uint16(vkCode),
	}

	log.Printf("🔍 Sending key press: %s (VK: 0x%02X)", key, vkCode)

	err := user32util.SendKeydbInput(keyInput, user32DLL)
	if err != nil {
		return fmt.Errorf("SendKeydbInput failed: %w", err)
	}

	return nil
}

func releaseKey(key string) error {
	vkCode, exists := keyMap[strings.ToLower(key)]
	if !exists {
		return fmt.Errorf("unknown key: %s", key)
	}

	keyInput := user32util.KeybdInput{
		WVK:     vkCode,
		DwFlags: 0x0002, // KEYEVENTF_KEYUP
	}

	log.Printf("🔍 Sending key release: %s (VK: 0x%02X)", key, vkCode)

	err := user32util.SendKeydbInput(keyInput, user32DLL)
	if err != nil {
		return fmt.Errorf("SendKeydbInput failed: %w", err)
	}

	return nil
}
