package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"own_wiki/src/ecv"
	e "own_wiki/src/events"
	t "own_wiki/src/platform/terminal"

	c "own_wiki/src/system/configuration"
	log "own_wiki/src/system/logger"
)

type MainView struct{}

func (mv *MainView) View(scene *ecv.Scene, yield func() bool) ecv.View {
	scene.CleanScreen()

	heading := ecv.NewHeading(1, "Titulo")
	scene.AddToScreen(heading)

	text := ecv.NewText("Hola")
	scene.AddToScreen(text)

	for i := range 5 * scene.FrameRate {
		if !yield() {
			fmt.Printf("Exiting at the first wait at %d\n", i)
			return nil
		}
	}

	text.ChangeText("Chau")

	for range 10 * scene.FrameRate {
		if !yield() {
			fmt.Println("Exiting at the last wait")
			return nil
		}
	}

	return nil
}

type TitleComponent struct {
	Title string
}

type TextComponent struct {
	Paragraphs []string
}

type FileEntity struct {
	TitleComponent
	TextComponent
}

/*
Con esto podemos definir 3 funciones, que fuerzan al usuario a establecer
todo lo que deberia hacer, es la API/contrato, que necesitan cumplir para que
el sistema funcione. Estas son
  - Funcion para ingresar los struct que corresponden como componentes
  - Funcion que recibe un archivo de texto, con toda la metadata, y el sistema
    para ingresar las entidades
  - Funcion para ingresar el par (entidad, []view)
*/
func SimulatedUser(ecvSystem *ecv.ECV) {
	// Registrar los componentes -> Esto se traduce en las estructuras de la base de datos
	ecvSystem.RegisterComponent(TitleComponent{})
	ecvSystem.RegisterComponent(TextComponent{})

	// Registrar las entidades
	/*
		Ahora que estan los componentes, podemos correr una funcion generada
		por el usuario que lea los archivos que tiene, y cree las entidades a
		partir de estos archivos.
	*/

	// Registrar las views
	mainView := &MainView{}
	ecvSystem.AssignCurrentView(mainView)
}

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

func Loop(config c.UserConfig, wg *sync.WaitGroup) {
	ecvSystem := ecv.NewECV(config) // Creamos event queue que va a ser un channel
	defer ecvSystem.Close()

	sigTermChannel := HandleSigTerm(ecvSystem.EventQueue)
	defer close(sigTermChannel)

	// Registrar estructura dadas por el usuario, y genera las views
	SimulatedUser(ecvSystem)

	platform := t.NewTerminal()
	defer platform.Close()

	wg.Add(1)
	go platform.HandleInput(ecvSystem.EventQueue, wg)

	ticker := time.NewTicker(time.Duration(1000/config.TargetFrameRate) * time.Millisecond)

	// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
	for range ticker.C {
		representation, ok := ecvSystem.GenerateFrame()
		if !ok {
			break
		}

		platform.Render(representation)
	}
}

func main() {
	if err := log.CreateLogger("logs/logger.txt", log.VERBOSE); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer log.Close()

	var waitGroup sync.WaitGroup
	Loop(c.UserConfig{
		TargetFrameRate: 1,
	}, &waitGroup)
	waitGroup.Wait()
}
