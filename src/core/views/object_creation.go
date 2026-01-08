package views

type ObjectCreator interface {
	NewScene(view View, worldConfig WorldConfiguration) *scene
	NewCamera() *camera

	NewLayout() *layout
}

type DefaultObjectCreator struct{}

func (DefaultObjectCreator) NewScene(view View, worldConfig WorldConfiguration) *scene {
	return newScene(view, worldConfig)
}

func (DefaultObjectCreator) NewCamera() *camera {
	return newCamera()
}

func (DefaultObjectCreator) NewLayout() *layout {
	return newLayout()
}
