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
		log.Debug("Waiting for Platform to handle all pending gorutings")
		wg.Done()
		log.Debug("Platform finished to handle input")
	}()

	ticker := time.NewTicker(time.Duration(1000/config.TargetFrameRate) * time.Millisecond)

	// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
	for range ticker.C {
		representation, ok := ecvSystem.GenerateFrame()
		if !ok {
			log.Debug("Leaving representation")
			break
		}

		platform.Render(representation)
	}

	log.Debug("Loop finished")
}
