package views

type WorldConfiguration struct {
}

func DefaultWorldConfiguration() WorldConfiguration {
	return WorldConfiguration{}
}

type World struct {
	MainCamera Camera

	configuration WorldConfiguration
}

func NewWorld(configuration WorldConfiguration) *World {
	return &World{}
}

func (w *World) GetConfiguration() WorldConfiguration {
	return w.configuration
}

func (w *World) Clear() {}

func (w *World) Render() SceneDescription {
	return nil
}
