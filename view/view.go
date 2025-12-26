package view

type SceneRepresentation []any

type View interface {
	View(scene *Scene, yield func() bool) View
}

type Heading struct {
	Level uint8
	Data  string
}

func NewHeading(level uint8, text string) *Heading {
	return &Heading{
		Level: level,
		Data:  text,
	}
}

type Text struct {
	Data string
}

func NewText(text string) *Text {
	return &Text{
		Data: text,
	}
}

func (t *Text) ChangeText(text string) {
	t.Data = text
}

type Link struct {
	data string
}
