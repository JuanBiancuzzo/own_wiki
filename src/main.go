package main

import (
	"time"

	e "own_wiki/ecv"
	htmx "own_wiki/platform/htmx"
)

func main() {
	ecv := e.NewECV() // Creamos event queue que va a ser un channel
	defer ecv.Close()

	platform := htmx.NewHTMX()
	defer platform.Close()

	go platform.HandleInput(ecv.EventQueue)

	// 60 Hz
	targetDuration := time.Duration(1000/60) * time.Millisecond
	ticker := time.NewTicker(targetDuration)

	// tick loop
	for range ticker.C {

	}
}
