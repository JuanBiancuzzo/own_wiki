package drawCommands

/*
For now we have as posible draw commands
 * DrawRectangle
 * DrawText
 * DrawPath

We could add as draw commands:
 * Images (using the path to the image)
 * Shaders
 * Polygons
 * 3D objects
*/
type DrawCommand isDrawCommand

type isDrawCommand interface {
	isDrawCommand()
}
