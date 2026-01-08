package views

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
)

type FnYield func() <-chan []e.Event

type View interface {
	View(world *World, creator ObjectCreator, yield FnYield) (nextView View)
}

type Renderable interface {
	Render() SceneDescription
}
