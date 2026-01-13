package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

// This structure is capable of waking the state machine define by the sequence
// of views
type ViewWalker[Data any] interface {
	// Se asume que para este punto la view ya fue preloadeada
	InitializeView(view View[Data])

	// Avanza la escena el siguiente frame
	WalkScene(events []e.Event)

	// Renderiza el frame
	Render() SceneRepresentation
}

type LocalWalker[Data any] struct {
	World *World
	Data  Data

	EventChannel chan []e.Event
	FrameChannel chan bool
	NextView     chan View[Data]
}

func NewLocalWalker[Data any](initialView View[Data], world *World, data Data) *LocalWalker[Data] {
	walker := &LocalWalker[Data]{
		World: world,
		Data:  data,
	}

	walker.InitializeView(initialView)
	return walker
}

func (lw *LocalWalker[Data]) InitializeView(view View[Data]) {
	lw.EventChannel = closeAndCreate(lw.EventChannel)
	lw.NextView = closeAndCreate(lw.NextView)
	lw.FrameChannel = closeAndCreate(lw.FrameChannel)

	var yield FnYield = func() <-chan []e.Event {
		lw.FrameChannel <- true
		return lw.EventChannel
	}

	view.Preload(lw.Data)
	go func(world *World, data Data, yield FnYield, nextView chan View[Data]) {
		nextView <- view.View(world, data, yield)
	}(lw.World, lw.Data, yield, lw.NextView)
}

func (lw *LocalWalker[Data]) WalkScene(events []e.Event) {
	lw.EventChannel <- events

	keepAdvancing := true

	for keepAdvancing {
		select {
		case <-lw.FrameChannel:
			keepAdvancing = false

		case view := <-lw.NextView:
			lw.InitializeView(view)
		}
	}
}

func (lw *LocalWalker[Data]) Render() SceneRepresentation {
	return lw.World.Render()
}

func closeAndCreate[T any](channel chan T) chan T {
	if channel != nil {
		close(channel)
	}
	return make(chan T)
}
