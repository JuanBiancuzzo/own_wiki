package scene

type Camera struct {
	PerspectiveMatrix [][]float32
	ScreenLayout      *Layout
}

func NewGeneralCamera(perspectiveMatrix [][]float32, layout *Layout) *Camera {
	return &Camera{
		PerspectiveMatrix: perspectiveMatrix,
		ScreenLayout:      layout,
	}
}

func New2DCamera(layout *Layout) *Camera {
	var identityMatrix = [][]float32{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	return NewGeneralCamera(identityMatrix, layout)
}
