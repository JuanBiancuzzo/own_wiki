package htmx

import (
	ecv "own_wiki/ecv"
	e "own_wiki/events"
	p "own_wiki/platform"
)

type HTMXPlatform struct{}

func NewHTMX() p.Platform {
	return &HTMXPlatform{}
}

func (hp *HTMXPlatform) HandleInput(eventQueue chan e.Event) {}

func (hp *HTMXPlatform) Render(view ecv.SceneRepresentation) {}

func (hp *HTMXPlatform) Close() {}
