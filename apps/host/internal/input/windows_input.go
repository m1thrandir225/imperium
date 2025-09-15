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
