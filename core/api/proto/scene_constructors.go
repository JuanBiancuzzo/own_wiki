package core

// ---+--- SceneDescription ---+---

func New2DScene(layout *Layout) *SceneDescription {
	return &SceneDescription{
		MainCamera: New2DCamera(layout),
	}
}

// ---+--- Camera ---+---
var IdentityMatrix = []float32{
	1, 0, 0, 0,
	0, 1, 0, 0,
	0, 0, 1, 0,
	0, 0, 0, 1,
}

func New2DCamera(layout *Layout) *Camera {
	return &Camera{
		Perspective: &Camera_PerspectiveMatrix{
			Value: IdentityMatrix,
		},
		ScreenLayout: layout,
	}
}

// ---+--- SceenLayout ---+---
func NewLayout(objects ...*SceneObject) *Layout {
	return &Layout{
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
