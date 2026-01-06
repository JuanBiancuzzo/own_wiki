package views

type Scene[Data any] struct {
	Walker ViewWalker[Data]
}

func NewScene[Data any](walker ViewWalker[Data]) *Scene[Data] {
	return &Scene[Data]{
		Walker: walker,
	}
}

func (s *Scene[_]) Render() SceneRepresentation {
	return s.Walker.Render()
}
