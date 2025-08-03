package fs

import (
	"fmt"
	"os"
	db "own_wiki/system_protocol/baseDeDatos"
	u "own_wiki/system_protocol/utilidades"
	"slices"
	"sync"
)

const CANTIDAD_WORKERS = 10

var DIRECTORIOS_IGNORAR = []string{".git", ".configuracion", ".github", ".obsidian", ".trash"}

type Root struct {
	Path     string
	Archivos map[string]*Archivo
}

func NewRoot(path string) *Root {
	return &Root{
		Archivos: make(map[string]*Archivo),
		Path:     path,
	}
}

func EstablecerDirectorio(dirOrigen string, infoArchivos *db.InfoArchivos, canalMensajes chan string) *Root {
	var waitArchivos sync.WaitGroup
	root := NewRoot(dirOrigen)

	canalInput := make(chan string, CANTIDAD_WORKERS)
	waitArchivos.Add(1)
	go DividirTrabajo(canalInput, root, infoArchivos, &waitArchivos, canalMensajes)

	colaDirectorios := u.NewCola[string]()
	colaDirectorios.Encolar("")
	canalMensajes <- fmt.Sprintf("El directorio para trabajar va a ser: %s\n", root.Path)

	for directorioPath := range colaDirectorios.DesencolarIterativamente {
		archivos, err := os.ReadDir(fmt.Sprintf("%s/%s", root.Path, directorioPath))
		if err != nil {
			canalMensajes <- fmt.Sprintf("Se tuvo un error al leer el directorio dando el error: %v", err)
			break
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

	close(canalInput)
	waitArchivos.Wait()
	return root
}

func DividirTrabajo(canalInput chan string, root *Root, info *db.InfoArchivos, wg *sync.WaitGroup, canalMensajes chan string) {
	var mutexRoot sync.Mutex

	procesarArchivo := func(path string) {
		if archivo, err := NewArchivo(root, path, info, canalMensajes); err != nil {
			canalMensajes <- fmt.Sprintf("Se tuvo un error al crear un archivo, con error: %v", err)
		} else {
			mutexRoot.Lock()
			root.Archivos[path] = archivo
			mutexRoot.Unlock()
		}
	}

	canalesInput := make([]chan string, CANTIDAD_WORKERS)
	var waitWorkers sync.WaitGroup

	waitWorkers.Add(CANTIDAD_WORKERS)
	for i := range CANTIDAD_WORKERS {
		canalesInput[i] = make(chan string, 5)
		worker := u.NewWorker(canalesInput[i], procesarArchivo, &waitWorkers)
		go worker.Ejecutar()
	}

	contador := 0
	for input := range canalInput {
		canalesInput[contador] <- input
		contador = (contador + 1) % CANTIDAD_WORKERS
	}

	for i := range CANTIDAD_WORKERS {
		close(canalesInput[i])
	}

	waitWorkers.Wait()
	wg.Done()
}

func (r *Root) EncontrarArchivo(path string) (*Archivo, error) {
	if archivo, ok := r.Archivos[path]; ok {
		return archivo, nil
	}
	return nil, fmt.Errorf("no existe el archivo con ese path")
}
