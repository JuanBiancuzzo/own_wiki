package ecv

import (
	e "own_wiki/events"
	v "own_wiki/view"
)

type ECV struct {
	EventQueue chan e.Event
}

func NewECV() *ECV {
	return &ECV{
		EventQueue: make(chan e.Event),
	}
}

func (ecv *ECV) GenerateFrame() (v.ViewRepresentation, bool) {
	return nil, true
}

func (ecv *ECV) Close() {}
