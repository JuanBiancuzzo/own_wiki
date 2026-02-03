package scene

type Camera struct {
	PerspectiveMatrix Mat4[float64]

	// The layout would be represented by a quad that is always on top of
	// everything else, independent from the camera.
	// In the layout could have another camera and it can have other objects
	ScreenLayout *Layout
	Objects      []Object
}

/*
We have custom constructors for:
  - Orthographic

We could create constructors for:
  - Isometric
  - Dimetric
  - Trimetric
  - Cabinet
  - Cavalier
  - Military
  - 1, 2, 3-points
  - Curvilinear
  - etc.
*/
func NewGeneralCamera(perspectiveMatrix Mat4[float64], layout *Layout) *Camera {
	return &Camera{
		PerspectiveMatrix: perspectiveMatrix,
		ScreenLayout:      layout,
		Objects:           []Object{layout},
	}
}

// ---+--- Create custom cameras ---+---

func NewOrthographicCamera(layout *Layout) *Camera {
	var identityMatrix = Mat4[float64]{
		A11: 1, A22: 1, A33: 1, A44: 1,
	}

	return NewGeneralCamera(identityMatrix, layout)
}

// ---+--- Funcionality ---+---

func (c *Camera) GenerateCameraDescription() []*CameraDescription {

	return []*CameraDescription{}
}
