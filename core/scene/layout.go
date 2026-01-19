package scene

type Layout struct {
	Objects []SceneObject
}

func NewLayout(objects ...SceneObject) *Layout {
	return &Layout{
		Objects: objects,
	}
}
