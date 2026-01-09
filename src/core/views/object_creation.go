package views

type ObjectCreator interface {
	NewScene(view View, worldConfig WorldConfiguration) *Scene
	NewCamera() *Camera

	NewLayout() *Layout
}

type DefaultObjectCreator struct{}

func (DefaultObjectCreator) NewScene(view View, worldConfig WorldConfiguration) *Scene {
	return DefaultNewScene(view, worldConfig)
}

func (DefaultObjectCreator) NewCamera() *Camera {
	return newCamera()
}

func (DefaultObjectCreator) NewLayout() *Layout {
	return newLayout()
}
