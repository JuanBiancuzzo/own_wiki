package drawCommands

type DrawRectangle struct {
	// Position
	// Size
	// Background color
	// rounded corners
}

func NewDrawRectangle() DrawRectangle {
	return DrawRectangle{}
}

func (dr DrawRectangle) isDrawCommand() {}
