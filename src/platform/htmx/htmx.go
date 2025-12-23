package htmx

import (
	"fmt"
	ecv "own_wiki/ecv"
	e "own_wiki/events"
	p "own_wiki/platform"
	"sync"
)

type HTMXPlatform struct{}

func NewHTMX() p.Platform {
	return &HTMXPlatform{}
}

func (hp *HTMXPlatform) HandleInput(eventQueue chan e.Event, wg *sync.WaitGroup) {
	wg.Done()
}

func (hp *HTMXPlatform) Render(view ecv.SceneRepresentation) {
	if len(view) > 0 {
		fmt.Println("Mostrando Screen Representation")
		for i, value := range view {
			fmt.Printf("%d: %v", i, value)
		}
		fmt.Println("")
	}
}

func (hp *HTMXPlatform) Close() {}
