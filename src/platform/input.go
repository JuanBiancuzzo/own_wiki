package platform

import (
	ecv "own_wiki/ecv"
	e "own_wiki/events"
)

type Platform interface {
	HandleInput(chan e.Event)

	Render(ecv.SceneRepresentation)

	Close()
}
