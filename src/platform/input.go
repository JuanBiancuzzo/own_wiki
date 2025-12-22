package platform

import (
	e "own_wiki/events"
)

type Platform interface {
	HandleInput(chan e.Event)

	Close()
}
