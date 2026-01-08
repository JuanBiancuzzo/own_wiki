package views

import e "github.com/JuanBiancuzzo/own_wiki/src/core/events"

type scene struct {
	InnerWorld *World
}

func newScene(view View, worldConfig WorldConfiguration) *scene {
	world := NewWorld(worldConfig)

	return &scene{
		InnerWorld: world,
	}
}

func (s scene) StepView(events []e.Event) bool {
	return false
}

func (s scene) Render() SceneDescription {
	return s.InnerWorld.Render()
}
