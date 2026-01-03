package views

type Scene struct {
	Walker ViewWalker
}

func NewScene(walker ViewWalker) *Scene {
	return &Scene{
		Walker: walker,
	}
}

func (s *Scene) Render() SceneRepresentation {
	return s.Walker.Render()
}
