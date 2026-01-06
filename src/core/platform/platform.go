package platform

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"
)

type Platform interface {
	HandleInput(chan e.Event)

	Render(v.SceneRepresentation)

	Close()
}
