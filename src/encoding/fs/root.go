package fs

import (
	"fmt"
	"os"
	t "own_wiki/system_protocol/tablas"
	u "own_wiki/system_protocol/utilidades"
	"slices"
	"sync"
)

const CANTIDAD_WORKERS = 15

var DIRECTORIOS_IGNORAR = []string{".git", ".configuracion", ".github", ".obsidian", ".trash"}

func RecorrerDirectorio(dirOrigen string, tablas *t.Tablas, canalMensajes chan string) error {
	var waitArchivos sync.WaitGroup

	canalInput := make(chan string, CANTIDAD_WORKERS)
	defer close(canalInput)
	waitArchivos.Add(1)
	defer waitArchivos.Wait()

	procesarArchivo := func(path string) {
		if err := CargarArchivo(dirOrigen, path, tablas, canalMensajes); err != nil {
			canalMensajes <- fmt.Sprintf("Se tuvo un error al crear un archivo en el path: '%s', con error: %v", path, err)
		}
	}
	go u.DividirTrabajo(canalInput, CANTIDAD_WORKERS, procesarArchivo, &waitArchivos)

	colaDirectorios := u.NewCola[string]()
	colaDirectorios.Encolar("")
	canalMensajes <- fmt.Sprintf("El directorio para trabajar va a ser: %s\n", dirOrigen)

	for directorioPath := range colaDirectorios.DesencolarIterativamente {
		archivos, err := os.ReadDir(fmt.Sprintf("%s/%s", dirOrigen, directorioPath))
		if err != nil {
			return fmt.Errorf("se tuvo un error al leer el directorio '%s' dando el error: %v", fmt.Sprintf("%s/%s", dirOrigen, directorioPath), err)
		}

		for _, archivo := range archivos {
			nombreArchivo := archivo.Name()
			archivoPath := nombreArchivo
			if directorioPath != "" {
				archivoPath = fmt.Sprintf("%s/%s", directorioPath, archivoPath)
			}

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, nombreArchivo) {
				colaDirectorios.Encolar(archivoPath)

			} else if !archivo.IsDir() {
				canalInput <- archivoPath
			}
		}
	}

	return nil
}
