package drawCommands

/*
For now we have as posible draw commands
 * DrawRectangle
 * DrawText

We could add as draw commands:
 * Path (a sequence of points like an SVG)
 * Images (using the path to the image)
 * Shaders
 * Polygons
 * 3D objects
*/
type DrawCommand isDrawCommand

type isDrawCommand interface {
	isDrawCommand()
}
