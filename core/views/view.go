package views

import (
	s "github.com/JuanBiancuzzo/own_wiki/core/scene"
)

type FnYield func() <-chan s.FrameInformation

type View[Data any] interface {
	// The view represents a continuous rendering of a scene, where each frame is
	// finish when the yield funtion is call. The Data any generics lets inyect
	// extra functionality, like a way to send events
	View(scene *s.Scene, data Data, yield FnYield) View[Data]
}
