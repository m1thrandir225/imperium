package input

type InputHandler interface {
	HandleCommand(cmd InputCommand)
}
