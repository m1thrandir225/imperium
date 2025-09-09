package input

type InputCommand struct {
	Type   string `json:"type"`   // "keyboard", "mouse"
	Action string `json:"action"` // "press", "release", "move", "click"
	Key    string `json:"key,omitempty"`
	X      int    `json:"x,omitempty"`
	Y      int    `json:"y,omitempty"`
	Button string `json:"button,omitempty"`
}
