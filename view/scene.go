package view

type Scene struct {
	Screen    []any
	FrameRate uint64
}

func NewScene(framRate uint64) *Scene {
	return &Scene{
		Screen:    []any{},
		FrameRate: framRate,
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
