package htmx

import (
	e "own_wiki/events"
	p "own_wiki/platform"
)

type HTMXPlatform struct{}

func NewHTMX() p.Platform {
	return &HTMXPlatform{}
}

func (hp *HTMXPlatform) HandleInput(chan e.Event) {}

func (hp *HTMXPlatform) Close() {}
