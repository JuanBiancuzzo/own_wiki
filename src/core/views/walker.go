package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

// This structure is capable of waking the state machine define by the sequence
// of views
type ViewWalker interface {
	// Se asume que para este punto la view ya fue preloadeada
	InitializeView(view View)

	// Preloadea la view en el caso de no haberlo sido
	Preload(uid ViewId, view View)

	// Avanza la escena el siguiente frame
	WalkScene(events []e.Event)

	// Renderiza el frame
	Render() SceneRepresentation
}

type LocalWalker struct {
	World        *World
	OutputEvents EventHandler
	RequestView  RequestView

	EventChannel chan []e.Event
	FrameChannel chan bool
	NextView     chan View

	Preloaded map[ViewId]View
}

func NewLocalWalker(initialView View, world *World, outputEvents EventHandler, requestView RequestView) *LocalWalker {
	walker := &LocalWalker{
		World:        world,
		OutputEvents: outputEvents,
		RequestView:  requestView,

		Preloaded: make(map[ViewId]View),
	}

	initialView.Preload(outputEvents)
	walker.InitializeView(initialView)

	return walker
}

func (lw *LocalWalker) InitializeView(view View) {
	lw.EventChannel = closeAndCreate(lw.EventChannel)
	lw.NextView = closeAndCreate(lw.NextView)
	lw.FrameChannel = closeAndCreate(lw.FrameChannel)

	var yield FnYield = func() []e.Event {
		lw.FrameChannel <- true
		return <-lw.EventChannel
	}

	go func(world *World, outputEvents EventHandler, yield FnYield, nextView chan View) {
		nextView <- view.View(world, outputEvents, yield)
	}(lw.World, lw.OutputEvents, yield, lw.NextView)
}

func (lw *LocalWalker) Preload(uid ViewId, view View) {
	if _, ok := lw.Preloaded[uid]; ok {
		return
	}

	view.Preload(lw.OutputEvents)
	lw.Preloaded[uid] = view
}

func (lw *LocalWalker) WalkScene(events []e.Event) {
	lw.EventChannel <- events

	keepAdvancing := true

	for keepAdvancing {
		select {
		case <-lw.FrameChannel:
			keepAdvancing = false

		case view := <-lw.NextView:
			uid, dataView := lw.RequestView.Request(view)
			if preloadedView, ok := lw.Preloaded[uid]; ok {
				dataView = preloadedView

				// Todo lo precargado se elimina ya que no fue elegido
				for vi := range lw.Preloaded {
					delete(lw.Preloaded, vi)
				}

			} else {
				dataView.Preload(lw.OutputEvents)
			}

			lw.InitializeView(dataView)
		}
	}
}

func (lw *LocalWalker) Render() SceneRepresentation {
	return lw.World.Render()
}

func closeAndCreate[T any](channel chan T) chan T {
	if channel != nil {
		close(channel)
	}
	return make(chan T)
}
