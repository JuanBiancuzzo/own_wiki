package core

// ---+--- SceneDescription ---+---
func NewGeneralScene(camera *CameraDescription) *SceneDescription {
	return &SceneDescription{
		MainCamera: camera,
	}
}

func New2DScene(layout *LayoutDescription) *SceneDescription {
	return &SceneDescription{
		MainCamera: New2DCamera(layout),
	}
}

// ---+--- Camera ---+---
func NewGeneralCamera(perspectiveMatrix []float32, layout *LayoutDescription) *CameraDescription {
	return &CameraDescription{
		Perspective: &CameraDescription_PerspectiveMatrix{
			Value: perspectiveMatrix,
		},
		ScreenLayout: layout,
	}
}

func New2DCamera(layout *LayoutDescription) *CameraDescription {
	var identityMatrix = []float32{
		1, 0, 0, 0,
		0, 1, 0, 0,
		0, 0, 1, 0,
		0, 0, 0, 1,
	}

	return NewGeneralCamera(identityMatrix, layout)
}

// ---+--- SceenLayout ---+---
func NewLayout(objects ...*SceneObject) *LayoutDescription {
	return &LayoutDescription{
		Objects: objects,
	}
}

// ---+--- SceenObject ---+---
func NewSceenTextObject(text string) *SceneObject {
	return &SceneObject{
		Object: &SceneObject_Text{
			Text: text,
		},
	}
}

func NewSceenImageObject(url, title string) *SceneObject {
	return &SceneObject{
		Object: &SceneObject_ImageInfo{
			ImageInfo: &SceneObject_Image{
				Url:   url,
				Title: title,
			},
		},
	}
}
