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
	if len(os.Args) <= 1 {
		fmt.Printf("No se ingresó que tipo de operacion se quiere ejecutar\n\tPuede ser -e para Encodear\n\tPuede ser -p para Procesar\n")
	}

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

	canalMensajes <- fmt.Sprint("Se tiene los argumentos sin limitar: ", os.Args)
	switch strings.TrimSpace(os.Args[1]) {
	case "-e":
		switch len(os.Args) {
		case 2:
			canalMensajes <- "No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de input"
		case 3:
			canalMensajes <- "No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de output"
		default:
			canalMensajes <- fmt.Sprintf("Se manda como input: %s, y manda como output: %s", os.Args[2], os.Args[3])
			en.Encodear(os.Args[2], os.Args[3], canalMensajes)
		}

	case "-p":
		switch len(os.Args) {
		case 2:
			canalMensajes <- "No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de input"
		default:
			canalMensajes <- fmt.Sprintf("Se manda como input: %s", os.Args[2])
			ej.Ejecutar(os.Args[2], canalMensajes)
		}

	default:
		canalMensajes <- fmt.Sprint("La eleccion elegida no fue una de las esperadas, esta fue: ", strings.TrimSpace(os.Args[1]))
	}

	close(canalMensajes)
	waitMensajes.Wait()

	fmt.Printf("Se termino el programa en: %s \n", time.Since(tiempoInicial))
}
