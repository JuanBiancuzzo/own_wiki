package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

type FnViewRequest func(View) View

type Scene struct {
	InnerWorld  *World
	ViewRequest FnViewRequest
	ObjectCreator

	eventQueue chan []e.Event
	nextView   chan View
	nextFrame  chan bool
}

func DefaultNewScene(view View, worldConfig WorldConfiguration) *Scene {
	scene := &Scene{
		InnerWorld:    NewWorld(worldConfig),
		ViewRequest:   func(view View) View { return view },
		ObjectCreator: DefaultObjectCreator{},
	}

	scene.initializeView(view)
	return scene
}

func NewCustomScene(view View, worldConfig WorldConfiguration, request FnViewRequest, objectCreator ObjectCreator) *Scene {
	scene := &Scene{
		InnerWorld:    NewWorld(worldConfig),
		ViewRequest:   request,
		ObjectCreator: DefaultObjectCreator{},
	}

	scene.initializeView(view)
	return scene
}

func (s *Scene) initializeView(view View) {
	view = s.ViewRequest(view)

	s.eventQueue = closeAndCreate(s.eventQueue)
	s.nextView = closeAndCreate(s.nextView)
	s.nextFrame = closeAndCreate(s.nextFrame)

	go func() {
		s.nextView <- view.View(s.InnerWorld, s.ObjectCreator, func() <-chan []e.Event {
			s.nextFrame <- true
			return s.eventQueue
		})
	}()
}

func (s *Scene) StepView(events []e.Event) bool {
	s.eventQueue <- events

	for {
		select {
		case <-s.nextFrame:
			return true

		case nextView := <-s.nextView:
			if nextView == nil {
				close(s.eventQueue)
				s.eventQueue = nil
				return false
			}

			s.initializeView(nextView)
		}
	}
}

func (s *Scene) Render() SceneDescription {
	return s.InnerWorld.Render()
}

func (s *Scene) Close() {
	if s.eventQueue != nil {
		close(s.eventQueue)
	}
}

func closeAndCreate[T any](channel chan T) chan T {
	if channel != nil {
		close(channel)
	}

	return make(chan T)
}
