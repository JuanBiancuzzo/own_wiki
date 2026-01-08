package views

type WorldConfiguration struct{}

func DefaultWorldConfiguration() WorldConfiguration {
	return WorldConfiguration{}
}

type World struct {
	MainCamera camera
}

func NewWorld(configuration WorldConfiguration) *World {
	return &World{}
}

func (w *World) Clear() {}

func (w *World) Render() SceneDescription {
	return nil
}
