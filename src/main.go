package main

import (
	"time"

	e "own_wiki/ecv"
	htmx "own_wiki/platform/htmx"
)

type MainView struct{}

func (mv *MainView) View(scene *e.Scene, yield func() bool) e.View {

	return mv
}

func SimulatedUser(ecv *e.ECV) {
	// Registrar los componentes -> Esto se traduce en las estructuras de la base de datos

	// Registrar las entidades

	// Registrar las views
	mainView := &MainView{}
	ecv.AssignCurrentView(mainView)
}

func main() {
	ecv := e.NewECV() // Creamos event queue que va a ser un channel
	defer ecv.Close()

	// Registrar estructura dadas por el usuario, y genera las views
	SimulatedUser(ecv)

	platform := htmx.NewHTMX()
	defer platform.Close()

	go platform.HandleInput(ecv.EventQueue)

	targetFrameRate := 60
	ticker := time.NewTicker(time.Duration(1000/targetFrameRate) * time.Millisecond)

	// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
	for range ticker.C {
		representation, ok := ecv.GenerateFrame()
		if !ok {
			break
		}

		platform.Render(representation)
	}
}
