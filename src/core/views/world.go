package views

type WorldConfiguration any

type World struct {
	MainCamera Camera
}

func NewWorld(configuration WorldConfiguration) *World {
	return &World{}
}

func (w *World) Clear() {}

func (w *World) Render() SceneRepresentation {
	return nil
}
