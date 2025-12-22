package ecv

type SceneRepresentation any

type Scene struct{}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) GetRepresentation() SceneRepresentation {
	return nil
}
