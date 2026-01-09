package views

// Crear un layout engine

type Layout struct{}

func newLayout() *Layout {
	return &Layout{}
}

func (l *Layout) Add(element Renderable) {}
