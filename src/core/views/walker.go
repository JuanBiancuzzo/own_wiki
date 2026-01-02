package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

// This structure is capable of waking the state machine define by the sequence
// of views
type ViewId uint64

type FnViewRequest func(requestedView View) (uid ViewId, dataView View)

type ViewWaker struct {
	World        *World
	OutputEvents EventHandler
	Request      FnViewRequest

	EventChannel chan []e.Event
	SceneChannel chan SceneRepresentation
	NextView     chan View

	Preloaded map[ViewId]View
}

func NewViewWaker(view View, world *World, outputEvents EventHandler, request FnViewRequest) *ViewWaker {
	view.Preload(outputEvents)

	viewWalker := &ViewWaker{
		World:        world,
		OutputEvents: outputEvents,
		Request:      request,

		Preloaded: make(map[ViewId]View),
	}

	viewWalker.initializeView(view)

	return viewWalker
}

func (vw *ViewWaker) Preload(uid ViewId, view View) {
	if _, ok := vw.Preloaded[uid]; ok {
		return
	}

	view.Preload(vw.OutputEvents)
	vw.Preloaded[uid] = view
}

func (vw *ViewWaker) WalkScene(events []e.Event) (newScene SceneRepresentation) {
	vw.EventChannel <- events

	keepAdvancing := true

	for keepAdvancing {
		select {
		case newScene = <-vw.SceneChannel:
			keepAdvancing = false

		case view := <-vw.NextView:
			uid, dataView := vw.Request(view)
			if preloadedView, ok := vw.Preloaded[uid]; ok {
				dataView = preloadedView

				// Todo lo precargado se elimina ya que no fue elegido
				for vi := range vw.Preloaded {
					delete(vw.Preloaded, vi)
				}

			} else {
				dataView.Preload(vw.OutputEvents)
			}

			vw.initializeView(dataView)
		}
	}

	return newScene
}

func (vw *ViewWaker) initializeView(view View) {
	vw.EventChannel = closeAndCreate(vw.EventChannel)
	vw.NextView = closeAndCreate(vw.NextView)
	vw.SceneChannel = closeAndCreate(vw.SceneChannel)

	var yield FnYield = func(scene SceneRepresentation) []e.Event {
		vw.SceneChannel <- scene
		return <-vw.EventChannel
	}

	go func(world *World, outputEvents EventHandler, yield FnYield, nextView chan View) {
		nextView <- view.View(world, outputEvents, yield)
	}(vw.World, vw.OutputEvents, yield, vw.NextView)
}

func closeAndCreate[T any](channel chan T) chan T {
	if channel != nil {
		close(channel)
	}
	return make(chan T)
}
