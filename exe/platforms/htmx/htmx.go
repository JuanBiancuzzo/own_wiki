package htmx

import (
	"fmt"

	e "github.com/JuanBiancuzzo/own_wiki/core/events"
	p "github.com/JuanBiancuzzo/own_wiki/core/platform"
	v "github.com/JuanBiancuzzo/own_wiki/view"
)

type HTMXPlatform struct{}

func NewHTMX() p.Platform {
	return &HTMXPlatform{}
}

func (hp *HTMXPlatform) HandleInput(eventQueue chan e.Event) {}

func (hp *HTMXPlatform) Render(view v.SceneRepresentation) {
	if len(view) > 0 {
		fmt.Println("Mostrando Screen Representation")
		for i, value := range view {
			fmt.Printf("%d: %v", i, value)
		}
		fmt.Println("")
	}
}

func (hp *HTMXPlatform) Close() {}
