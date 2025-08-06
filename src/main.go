package main

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	ej "own_wiki/ejecucion"
	en "own_wiki/encoding"
)

func TerminarEjecucion(tiempoInicial time.Time, canalMensajes chan string, wg *sync.WaitGroup) {
}

func main() {
	tiempoInicial := time.Now()
	var waitMensajes sync.WaitGroup
	canalMensajes := make(chan string, 100)

	waitMensajes.Add(1)
	go func(canal chan string, wg *sync.WaitGroup) {
		for mensaje := range canal {
			fmt.Println(strings.TrimSpace(mensaje))
		}
		wg.Done()
	}(canalMensajes, &waitMensajes)

	if len(os.Args) <= 1 {
		ej.Ejecutar(canalMensajes)
	} else {
		switch len(os.Args) {
		case 2:
			canalMensajes <- "No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de input"
		default:
			en.Encodear(os.Args[2], canalMensajes)
		}
	}

	close(canalMensajes)
	waitMensajes.Wait()

	fmt.Printf("Se termino el programa en: %s \n", time.Since(tiempoInicial))
}
