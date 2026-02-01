package scene

type Camera struct {
	PerspectiveMatrix Mat4[float64]
	ScreenLayout      *Layout
}

func NewGeneralCamera(perspectiveMatrix Mat4[float64], layout *Layout) *Camera {
	return &Camera{
		PerspectiveMatrix: perspectiveMatrix,
		ScreenLayout:      layout,
	}
}

func New2DCamera(layout *Layout) *Camera {
	var identityMatrix = Matrix4x4[float64]{
		A11: 1, A22: 1, A33: 1, A44: 1,
	}

	return NewGeneralCamera(Mat4[float64](identityMatrix), layout)
}
