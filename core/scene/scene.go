package scene

type SceneDescription struct {
	MainCamera *Camera
}

func New2dScene(layout *Layout) *SceneDescription {
	return &SceneDescription{
		MainCamera: New2DCamera(layout),
	}
}

var IdentityMatrix = [][]float32{
	{1, 0, 0, 0},
	{0, 1, 0, 0},
	{0, 0, 1, 0},
	{0, 0, 0, 1},
}

type Camera struct {
	PerpectiveMatrix [][]float32
	ScreenLayout     *Layout
}

func New2DCamera(layout *Layout) *Camera {
	return &Camera{
		PerpectiveMatrix: IdentityMatrix,
		ScreenLayout:     layout,
	}
}

type Layout struct {
	Objects []ScreenObject
}

func NewLayout(objects ...ScreenObject) *Layout {
	return &Layout{
		Objects: objects,
	}
}

type ScreenObject isScreenObject

type isScreenObject interface {
	isScreenObject()
}

type TextObject struct {
	Text string
}

func (*TextObject) isScreenObject() {}

func NewTextObject(text string) *TextObject {
	return &TextObject{
		Text: text,
	}
}

type ImageObject struct {
	Url   string
	Title string
}

func (*ImageObject) isScreenObject() {}

func NewImageObject(url, title string) *ImageObject {
	return &ImageObject{
		Url:   url,
		Title: title,
	}
}
