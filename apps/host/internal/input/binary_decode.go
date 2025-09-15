package input

import "encoding/binary"

const (
	inpTypeKeyboard    = 0
	inpTypeMouseMove   = 1
	inpTypeMouseButton = 2
	inpTypeWheel       = 3

	actPress   = 0
	actRelease = 1
	actMove    = 2
	actClick   = 3

	btnNone   = 0
	btnLeft   = 1
	btnRight  = 2
	btnMiddle = 3
)

//DecodeInputCommand decodes a binary input command to the InputCommand struct
func DecodeInputCommand(b []byte) (InputCommand, bool) {
	if len(b) < 10 {
		return InputCommand{}, false
	}

	t := b[0]
	a := b[1]
	btn := b[2]

	key := binary.LittleEndian.Uint16(b[4:6])
	x := int16(binary.LittleEndian.Uint16(b[6:8]))
	y := int16(binary.LittleEndian.Uint16(b[8:10]))

	var cmd InputCommand

	switch t {
	case inpTypeKeyboard:
		if a == actPress {
			cmd.Action = "press"
		} else {
			cmd.Action = "release"
		}

		cmd.Key = keyCodeToString(key)
	case inpTypeMouseMove:
		cmd.Action = "move"
		cmd.X = int(x)
		cmd.Y = int(y)
	case inpTypeMouseButton:
		cmd.Type = "mouse"
		if a == actPress {
			cmd.Action = "press"
		} else {
			cmd.Action = "release"
		}
		cmd.Button = buttonToString(btn)
	case inpTypeWheel:
		cmd.Type = "mouse"
		cmd.Action = "scroll"
		cmd.Y = int(y)
	default:
		return InputCommand{}, false
	}

	return cmd, true
}

func buttonToString(b byte) string {
	switch b {
	case btnLeft:
		return "left"
	case btnRight:
		return "right"
	case btnMiddle:
		return "middle"
	default:
		return "none"
	}
}

func keyCodeToString(code uint16) string {
	return vkToName[code]
}
