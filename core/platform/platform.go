package platform

import (
	"github.com/JuanBiancuzzo/own_wiki/core/ecv"
	e "github.com/JuanBiancuzzo/own_wiki/core/events"
)

type Platform interface {
	HandleInput(chan e.Event)

	Render(ecv.SceneRepresentation)

	Close()
}
