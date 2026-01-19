package scene

// The types that holds the interface
// 	*TextObject
// 	*ImageObject
type SceneObject isScreenObject

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
