package platform

import (
	e "own_wiki/events"
	v "own_wiki/view"
)

type Platform interface {
	HandleInput(chan e.Event)

	Render(v.ViewRepresentation)

	Close()
}
