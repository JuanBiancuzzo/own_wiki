package htmx

import (
	e "own_wiki/events"
	p "own_wiki/platform"
	v "own_wiki/view"
)

type HTMXPlatform struct{}

func NewHTMX() p.Platform {
	return &HTMXPlatform{}
}

func (hp *HTMXPlatform) HandleInput(eventQueue chan e.Event) {}

func (hp *HTMXPlatform) Render(view v.ViewRepresentation) {}

func (hp *HTMXPlatform) Close() {}
