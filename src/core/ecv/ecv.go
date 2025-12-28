package ecv

import (
	e "github.com/JuanBiancuzzo/own_wiki/src/core/views/events"
)

type EventHandler interface {
	// This method is to be capable to send modifications to the system in a view
	PushEvent(event e.Event) error
}

type ECV any
