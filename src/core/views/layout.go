package views

// Crear un layout engine

type layout struct{}

func newLayout() *layout {
	return &layout{}
}

func (l *layout) Add(element Renderable) {}
