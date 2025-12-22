package ecv

import (
	e "own_wiki/events"
)

type ECV struct {
	EventQueue chan e.Event
}

func NewECV() *ECV {
	return &ECV{
		EventQueue: make(chan e.Event),
	}
}

func (ecv *ECV) Close() {}
