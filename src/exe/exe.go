package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "embed"

	ps "github.com/JuanBiancuzzo/own_wiki/src/exe/platforms"
	vs "github.com/JuanBiancuzzo/own_wiki/src/exe/views"

	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	p "github.com/JuanBiancuzzo/own_wiki/src/core/platform"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

	c "github.com/JuanBiancuzzo/own_wiki/src/core/systems/configuration"
	log "github.com/JuanBiancuzzo/own_wiki/src/core/systems/logging"
)

// Cambiarlo a argumento
const USER_CONFIG_PATH = ""

func HandleSigTerm(eventQueue chan e.Event) chan os.Signal {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		if _, ok := <-signals; ok {
			log.Debug("Getting ctrl+c interrupt")
			eventQueue <- e.NewCloseEvent("Ctrl+c interrupt")
		}
	}()

	return signals
}

func main() {
	eventQueue := make(chan e.Event)
	sigTermChannel := HandleSigTerm(eventQueue)
	defer close(sigTermChannel)

	if err := c.LoadUserConfiguration(USER_CONFIG_PATH); err != nil {
		log.Error("Failed to load user configuration")
		return
	}

	var platform p.Platform
	switch "Terminal" {
	case "Terminal":
		platform = ps.NewTerminalPlatform()
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		platform.HandleInput(eventQueue)
		wg.Done()
		log.Debug("Platform finished to handle input")
	}()

	eventsPerFrame := make(chan []e.Event)
	defer close(eventsPerFrame)

	go func() {
		// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
		ticker := time.NewTicker(time.Duration(1000/c.UserConfig.TargetFrameRate) * time.Millisecond)

		acumulatedEvents := []e.Event{}
		keepReading := true

		for keepReading {
			select {
			case events, ok := <-eventQueue:
				if !ok {
					keepReading = false
					break
				}
				acumulatedEvents = append(acumulatedEvents, events)

			case <-ticker.C:
				eventsPerFrame <- []e.Event{}

				// Simula eliminar los eventos utilizados
				acumulatedEvents = []e.Event{}
			}
		}
		close(eventQueue)
		log.Debug("Events handler for scene close")
	}()

	var yield v.FnYield = func() <-chan []e.Event {
		return eventsPerFrame
	}

	mainView := vs.NewMainView()
	mainView.View(v.NewWorld(v.WorldConfiguration(0)), yield)
	log.Info("Finish the main loop, closing app")
}
