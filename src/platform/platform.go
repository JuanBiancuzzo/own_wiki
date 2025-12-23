package platform

import (
	ecv "own_wiki/ecv"
	e "own_wiki/events"
	"sync"
)

type Platform interface {
	HandleInput(chan e.Event, *sync.WaitGroup)

	Render(ecv.SceneRepresentation)

	Close()
}
