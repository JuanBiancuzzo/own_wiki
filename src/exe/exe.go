package main

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "embed"

	// p "github.com/JuanBiancuzzo/own_wiki/src/exe/platforms"
	vs "github.com/JuanBiancuzzo/own_wiki/src/exe/views"

	"github.com/JuanBiancuzzo/own_wiki/src/core/api"
	"github.com/JuanBiancuzzo/own_wiki/src/core/ecv"
	e "github.com/JuanBiancuzzo/own_wiki/src/core/events"
	u "github.com/JuanBiancuzzo/own_wiki/src/core/user"
	v "github.com/JuanBiancuzzo/own_wiki/src/core/views"

	log "github.com/JuanBiancuzzo/own_wiki/src/core/systems/logging"
)

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

	// Registrar estructura dadas por el usuario, y genera las views
	userStructureData, err := u.GetUserDefineData(config.UserDefineDataDirectory)
	if err != nil {
		log.Error("Failed to get user define data plugin, with error: '%v'", err)
		return
	}
	defer userStructureData.Close()

	var ecv *ecv.ECV
	if ecv, err = userStructureData.RegisterStructures(); err != nil {
		log.Error("Failed to get components from user, with error: '%v'", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		platform.HandleInput(eventQueue)
		wg.Done()
		log.Debug("Platform finished to handle input")
	}()

	wg.Add(1)
	eventsPerFrame := make(chan []e.Event)
	defer close(eventsPerFrame)

	go func() {
		// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
		ticker := time.NewTicker(time.Duration(1000/config.TargetFrameRate) * time.Millisecond)

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
		wg.Done()
		log.Debug("Events handler for scene close")
	}()

	mainView := &vs.MainView{
		UserView: vs.NewUserPluginWalker(userStructureData.Plugin, ecv),
	}
	mainWalker := v.NewLocalWalker[api.OWData](mainView, v.NewWorld(v.WorldConfiguration(0)), ecv)

	for events := range eventsPerFrame {
		if !mainWalker.WalkScene(events) {
			break
		}

		scene := mainWalker.Render()
		platform.Render(scene)
	}

	log.Debug("Loop finished")
}
