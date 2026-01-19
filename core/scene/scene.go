package scene

type Scene struct {
	MainCamera *Camera
}

func NewGeneralScene(camera *Camera) *Scene {
	return &Scene{MainCamera: camera}
}

func New2DScene(layout *Layout) *Scene {
	return &Scene{
		MainCamera: New2DCamera(layout),
	}
}
