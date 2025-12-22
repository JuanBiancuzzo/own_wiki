package main

import (
	"time"

	e "own_wiki/ecv"
	htmx "own_wiki/platform/htmx"
)

func main() {
	ecv := e.NewECV() // Creamos event queue que va a ser un channel
	defer ecv.Close()

	// Registrar estructura dadas por el usuario, y genera las views

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
