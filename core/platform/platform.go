package platform

import (
	e "github.com/JuanBiancuzzo/own_wiki/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/view"
)

type Platform interface {
	HandleInput(chan e.Event)

	Render(v.SceneRepresentation)

	Close()
}
