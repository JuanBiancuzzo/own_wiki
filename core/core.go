package core

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/JuanBiancuzzo/own_wiki/core/ecv"
	e "github.com/JuanBiancuzzo/own_wiki/core/events"
	p "github.com/JuanBiancuzzo/own_wiki/core/platform"
	c "github.com/JuanBiancuzzo/own_wiki/core/system/configuration"
	log "github.com/JuanBiancuzzo/own_wiki/core/system/logger"
	u "github.com/JuanBiancuzzo/own_wiki/core/user"

	v "github.com/JuanBiancuzzo/own_wiki/view"
)

func handleSigTerm(eventQueue chan e.Event) chan os.Signal {
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

func Loop(config c.UserConfig, platform p.Platform, wg *sync.WaitGroup) {
	ecvSystem := ecv.NewECV(config) // Creamos event queue que va a ser un channel
	defer ecvSystem.Close()

	sigTermChannel := handleSigTerm(ecvSystem.EventQueue)
	defer close(sigTermChannel)

	// Registrar estructura dadas por el usuario, y genera las views
	userDefineData, err := u.GetUserDefineData(config.UserDefineDataDirectory)
	if err != nil {
		log.Error("Failed to get user define data plugin, with error: '%v'", err)
		return
	}
	defer userDefineData.Close()

	if componentTypes, err := userDefineData.RegisterComponents(); err != nil {
		log.Error("Failed to get components from user, with error: '%v'", err)
		return

	} else {
		for _, component := range componentTypes {
			log.Info("Component register by the user is: %+v", *component)
		}
	}

	wg.Add(1)
	go func() {
		platform.HandleInput(ecvSystem.EventQueue)
		wg.Done()
		log.Debug("Platform finished to handle input")
	}()

	// Unificar esto
	eventChannel := make(chan []v.Event)
	wg.Add(1)
	go func() {
		// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
		ticker := time.NewTicker(time.Duration(1000/config.TargetFrameRate) * time.Millisecond)

		acumulatedEvents := []e.Event{}
		keepReading := true

		for keepReading {
			select {
			case events, ok := <-ecvSystem.EventQueue:
				if !ok {
					keepReading = false
					break
				}
				acumulatedEvents = append(acumulatedEvents, events)

			case <-ticker.C:
				eventChannel <- []v.Event{}

				// Simula eliminar los eventos utilizados
				acumulatedEvents = []e.Event{}
			}
		}

		close(eventChannel)
		wg.Done()
		log.Debug("Events handler for scene close")
	}()

	userDefineData.CreateView("MainView", ecvSystem.Scene)
	for events := range eventChannel {
		ops, err := userDefineData.AvanzarView(events)
		if err != nil {
			log.Debug("Failed to advance the view, with error: %v", err)
			break
		}

		platform.Render(*ops.SceneCaracteristics)

		if ops.EndScene == true {
			log.Debug("Leaving representation")
			break
		}

		if ops.ChangeScene != nil {
			userDefineData.CreateView(ops.ChangeScene.ViewName, ecvSystem.Scene)
		}
	}

	log.Debug("Loop finished")
}
