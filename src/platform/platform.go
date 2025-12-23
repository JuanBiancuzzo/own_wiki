package platform

import (
	ecv "own_wiki/src/ecv"
	e "own_wiki/src/events"
	"sync"
)

type Platform interface {
	HandleInput(chan e.Event, *sync.WaitGroup)

	Render(ecv.SceneRepresentation)

	Close()
}
