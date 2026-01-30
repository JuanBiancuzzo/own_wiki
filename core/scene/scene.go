package scene

// We should contamplate the idea of having a way to represent the scene
// instead of sharing all the scene itself
type SceneCtx struct {
	MainCamera *Camera
}

func NewGeneralScene(camera *Camera) *SceneCtx {
	return &SceneCtx{MainCamera: camera}
}

func New2DScene(layout *Layout) *SceneCtx {
	return &SceneCtx{
		MainCamera: New2DCamera(layout),
	}
}
