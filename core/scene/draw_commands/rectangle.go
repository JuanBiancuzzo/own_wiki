package drawCommands

type DrawRectangle struct {
	// Position
	// Size
	// Background color
	// rounded corners // for now we dont include it
}

func NewDrawRectangle() DrawRectangle {
	return DrawRectangle{}
}

func (dr DrawRectangle) isDrawCommand() {}
