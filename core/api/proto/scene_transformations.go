package core

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

func ConvertFromSystemScene(scene s.Scene) (sceneDescription *SceneDescription, err error) {
	return sceneDescription, err
}

func (so *SceneObject) ConvertToSystemObject() (object s.SceneObject, err error) {
	switch so.GetObject().(type) {
	case *SceneObject_Text:
		text := so.GetText()
		object = s.NewTextObject(text)

	case *SceneObject_ImageInfo:
		imageInfo := so.GetImageInfo()
		object = s.NewImageObject(imageInfo.Url, imageInfo.Title)

	default:
		err = fmt.Errorf("Object type not a SceneObject, was of: %T", so.GetObject())
	}

	return object, err
}

func (ld *LayoutDescription) ConvertToSystemLayout() (*s.Layout, error) {
	objects := make([]s.SceneObject, len(ld.Objects))
	for i, objectDescription := range ld.Objects {
		if object, err := objectDescription.ConvertToSystemObject(); err != nil {
			return nil, fmt.Errorf("Failed to convert SceneObject (%s), with error: %v", objectDescription.String(), err)

		} else {
			objects[i] = object
		}
	}
	return s.NewLayout(objects...), nil
}

func (cd *CameraDescription) ConvertToSystemScene() (*s.Camera, error) {
	// The matrix should be of 4x4, then the length of the matrix should be 16
	matrixArray := cd.Perspective.Value
	if len(matrixArray) != 16 {
		return nil, fmt.Errorf("The perspective matrix is not of 4x4")
	}

	perspectiveMatrix := make([][]float32, 4)
	for i := range 4 {
		perspectiveMatrix[i] = matrixArray[i*4 : (i+1)*4]
	}

	if layout, err := cd.ScreenLayout.ConvertToSystemLayout(); err != nil {
		return nil, fmt.Errorf("Failed to convert ScreenLayout, with error: %v", err)

	} else {
		return s.NewGeneralCamera(perspectiveMatrix, layout), nil
	}
}

func (sd *SceneDescription) ConvertToSystemScene() (*s.Scene, error) {
	if camera, err := sd.MainCamera.ConvertToSystemScene(); err != nil {
		return nil, err

	} else {
		return s.NewGeneralScene(camera), nil
	}
}
