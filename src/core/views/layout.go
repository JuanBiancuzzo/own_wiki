package views

// Crear un layout engine

type Layout struct{}

func NewLayout() *Layout {
	return &Layout{}
}

func (l *Layout) Add(element Renderable) {}
