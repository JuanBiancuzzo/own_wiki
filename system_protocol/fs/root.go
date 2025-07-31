package fs

import (
	"fmt"
	"os"
	"own_wiki/system_protocol/db"
	ls "own_wiki/system_protocol/listas"
	"slices"
	"sync"
)

type Root struct {
	DirectorioRoot *Directorio
	Archivos       map[string]*Archivo
}

func EstablecerDirectorio(dirOrigen string, infoArchivos *db.InfoArchivos, canal chan string) *Root {
	var waitArchivos sync.WaitGroup

	directorioRoot := NewRoot(dirOrigen)
	root := &Root{
		DirectorioRoot: directorioRoot,
		Archivos:       make(map[string]*Archivo),
	}

	colaDirectorios := ls.NewCola[*Directorio]()
	colaDirectorios.Encolar(directorioRoot)
	canal <- fmt.Sprintf("El directorio para trabajar va a ser: %s\n", directorioRoot.Path)

	for !colaDirectorios.Vacia() {
		directorio, err := colaDirectorios.Desencolar()
		if err != nil {
			canal <- fmt.Sprintf("Se tuvo un error al operar sobre la queue con el error: %v", err)
			break
		}

		archivos, err := os.ReadDir(directorio.Path)
		if err != nil {
			canal <- fmt.Sprintf("Se tuvo un error al leer el directorio dando el error: %v", err)
			break
		}

		listaArchivos := ls.NewLista[string]()
		for _, archivo := range archivos {
			nombreArchivo := archivo.Name()
			archivoPath := fmt.Sprintf("%s/%s", directorio.Path, nombreArchivo)

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, nombreArchivo) {
				subDirectorio := NewDirectorio(directorio, archivoPath)
				directorio.AgregarSubdirectorio(nombreArchivo, subDirectorio)
				colaDirectorios.Encolar(subDirectorio)

			} else if !archivo.IsDir() {
				listaArchivos.Push(archivoPath)
			}
		}

		waitArchivos.Add(1)
		go func(lista *ls.Lista[string], root *Root, directorio *Directorio, infoArchivo *db.InfoArchivos, wg *sync.WaitGroup) {
			for archivoPath := range lista.Iterar {
				if archivo, err := NewArchivo(root, archivoPath, infoArchivos, canal); err != nil {
					canal <- fmt.Sprintf("Se tuvo un error al crear un archivo, con error: %v", err)

				} else {
					directorio.AgregarArchivo(archivo)
				}
			}
			wg.Done()
		}(listaArchivos, root, directorio, infoArchivos, &waitArchivos)
	}

	waitArchivos.Wait()
	directorioRoot.RelativizarPath(fmt.Sprintf("%s/", dirOrigen))

	for archivo := range directorioRoot.IterarArchivos {
		root.Archivos[archivo.Path] = archivo
	}

	return root
}

func (r *Root) EncontrarArchivo(path string) (*Archivo, error) {
	if archivo, ok := r.Archivos[path]; !ok {
		return nil, fmt.Errorf("no existe el archivo con ese path")
	} else {
		return archivo, nil
	}
}
