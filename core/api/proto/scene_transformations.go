package core

import (
	"fmt"

	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

func ConvertFromSystemObject(object s.SceneObject) (objectDescription *SceneObject, err error) {
	switch value := object.(type) {
	case *s.TextObject:
		objectDescription = NewSceenTextObject(value.Text)

	case *s.ImageObject:
		objectDescription = NewSceenImageObject(value.Url, value.Title)

	default:
		err = fmt.Errorf("Object type not a SceneObject, was of: %T", object)
	}

	return objectDescription, err
}

func ConvertFromSystemLayout(layout *s.Layout) (*LayoutDescription, error) {
	objectDescriptions := make([]*SceneObject, len(layout.Objects))
	for i, object := range layout.Objects {
		var err error
		if objectDescriptions[i], err = ConvertFromSystemObject(object); err != nil {
			return nil, fmt.Errorf("Failed to conver from object, with error: %v", err)
		}
	}
	return NewLayout(objectDescriptions...), nil
}

func ConvertFromSystemCamera(camera *s.Camera) (*CameraDescription, error) {
	perspectiveMatrix := make([]float32, 16)
	for i := range 4 {
		for j := range 4 {
			perspectiveMatrix[i+4*j] = camera.PerspectiveMatrix[i][j]
		}
	}

	if layout, err := ConvertFromSystemLayout(camera.ScreenLayout); err != nil {
		return nil, fmt.Errorf("Failed to convert from Layout, with error: %v", err)

	} else {
		return NewGeneralCamera(perspectiveMatrix, layout), nil
	}
}

func ConvertFromSystemScene(scene *s.Scene) (*SceneDescription, error) {
	if mainCamera, err := ConvertFromSystemCamera(scene.MainCamera); err != nil {
		return nil, fmt.Errorf("Failed to convert from Camera, with error: %v", err)

	} else {
		return NewGeneralScene(mainCamera), nil
	}
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
		var err error
		if objects[i], err = objectDescription.ConvertToSystemObject(); err != nil {
			return nil, fmt.Errorf("Failed to convert SceneObject (%s), with error: %v", objectDescription.String(), err)
		}
	}
	return s.NewLayout(objects...), nil
}

func (cd *CameraDescription) ConvertToSystemCamera() (*s.Camera, error) {
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
	if camera, err := sd.MainCamera.ConvertToSystemCamera(); err != nil {
		return nil, err

	} else {
		return s.NewGeneralScene(camera), nil
	}
}
