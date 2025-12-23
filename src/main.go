package main

import (
	"fmt"
	"sync"
	"time"

	e "own_wiki/ecv"
	t "own_wiki/platform/terminal"
	log "own_wiki/system/logger"
)

type MainView struct{}

func (mv *MainView) View(scene *e.Scene, yield func() bool) e.View {
	scene.CleanScreen()

	heading := e.NewText("Titulo")
	scene.AddToScreen(heading)

	text := e.NewText("Hola")
	scene.AddToScreen(text)

	for i := range 3 * 60 {
		if !yield() {
			fmt.Printf("Exiting at the first wait at %d\n", i)
			return nil
		}
	}

	text.ChangeText("Chau")

	for range 3 * 60 {
		if !yield() {
			fmt.Println("Exiting at the last wait")
			return nil
		}
	}

	return mv
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
func SimulatedUser(ecv *e.ECV) {
	// Registrar los componentes -> Esto se traduce en las estructuras de la base de datos
	ecv.RegisterComponent(TitleComponent{})
	ecv.RegisterComponent(TextComponent{})

	// Registrar las entidades
	/*
		Ahora que estan los componentes, podemos correr una funcion generada
		por el usuario que lea los archivos que tiene, y cree las entidades a
		partir de estos archivos.
	*/

	// Registrar las views
	mainView := &MainView{}
	ecv.AssignCurrentView(mainView)
}

func main() {
	if err := log.CreateLogger("logs/logger.txt", log.VERBOSE); err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer log.Close()

	var waitGroup sync.WaitGroup

	ecv := e.NewECV() // Creamos event queue que va a ser un channel
	defer ecv.Close()

	// Registrar estructura dadas por el usuario, y genera las views
	SimulatedUser(ecv)

	platform := t.NewTerminal()
	defer platform.Close()

	waitGroup.Add(1)
	go platform.HandleInput(ecv.EventQueue, &waitGroup)

	targetFrameRate := 60
	ticker := time.NewTicker(time.Duration(1000/targetFrameRate) * time.Millisecond)

	// Esto fuerza a que cada iteración como mínimo dure 1/FrameRate
	for range ticker.C {
		representation, ok := ecv.GenerateFrame()
		if !ok {
			fmt.Println("Not ok")
			break
		}

		platform.Render(representation)
	}

	waitGroup.Wait()
}
