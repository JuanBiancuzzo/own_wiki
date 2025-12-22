package ecv

type View interface {
	View(scene *Scene, yield func() bool) View
}
