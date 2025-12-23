package ecv

type Scene struct {
	Screen []any
}

func NewScene() *Scene {
	return &Scene{
		Screen: []any{},
	}
}

func (s *Scene) CleanScreen() {
	s.Screen = []any{}
}

func (s *Scene) AddToScreen(element any) {
	s.Screen = append(s.Screen, element)
}

func (s *Scene) GetRepresentation() SceneRepresentation {
	return s.Screen
}
