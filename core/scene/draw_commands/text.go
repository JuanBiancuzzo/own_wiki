package drawCommands

type DrawText struct {
}

func NewDrawText() DrawText {
	return DrawText{}
}

func (dr DrawText) isDrawCommand() {}
