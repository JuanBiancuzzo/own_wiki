package fs

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
	ls "own_wiki/system_protocol/listas"

	_ "github.com/go-sql-driver/mysql"
)

var DIRECTORIOS_IGNORAR = []string{".git", ".configuracion", ".github", ".obsidian", ".trash"}

type Directorio struct {
	Path           string
	Subdirectorios *ls.Lista[*Directorio]
	Archivos       *ls.Lista[*Archivo]
}

func NewDirectorio(path string) *Directorio {
	return &Directorio{
		Path:           path,
		Subdirectorios: ls.NewLista[*Directorio](),
		Archivos:       ls.NewLista[*Archivo](),
	}
}

func EstablecerDirectorio(root string, canal chan string) *Directorio {
	directorioRoot := NewDirectorio(root)

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

		for _, archivo := range archivos {
			archivoPath := fmt.Sprintf("%s/%s", directorio.Path, archivo.Name())

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, archivo.Name()) {
				nuevoDirectorio := NewDirectorio(archivoPath)
				directorio.AgregarSubdirectorio(nuevoDirectorio)
				colaDirectorios.Encolar(nuevoDirectorio)

			} else if !archivo.IsDir() {
				directorio.AgregarArchivo(NewArchivo(archivoPath))
			}
		}
	}

	return directorioRoot
}

func (d *Directorio) RelativizarPath(path string) {
	d.Path = strings.Replace(d.Path, path, "", 1)

	for archivo := range d.Archivos.Iterar {
		archivo.RelativizarPath(path)
	}

	for directorio := range d.Subdirectorios.Iterar {
		directorio.RelativizarPath(path)
	}
}

func (d *Directorio) AgregarSubdirectorio(directorio *Directorio) {
	d.Subdirectorios.Push(directorio)
}

func (d *Directorio) AgregarArchivo(archivo *Archivo) {
	d.Archivos.Push(archivo)
}

func (d *Directorio) ProcesarArchivos(wg *sync.WaitGroup, infoArchivos *db.InfoArchivos, canal chan string) {
	for subdirectorio := range d.Subdirectorios.Iterar {
		if subdirectorio.Vacio() {
			continue
		}

		wg.Add(1)
		go func(directorio *Directorio, wg *sync.WaitGroup) {
			directorio.ProcesarArchivos(wg, infoArchivos, canal)
			wg.Done()
		}(subdirectorio, wg)
	}

	for archivo := range d.Archivos.Iterar {
		archivo.Interprestarse(infoArchivos, canal)
	}
}

func (d *Directorio) Vacio() bool {
	return d.Subdirectorios.Vacia() && d.Archivos.Vacia()
}

func (d *Directorio) IterarArchivos(yield func(*Archivo) bool) {
	for archivo := range d.Archivos.Iterar {
		if !yield(archivo) {
			return
		}
	}

	for subdirectorio := range d.Subdirectorios.Iterar {
		subdirectorio.IterarArchivos(yield)
	}
}

func (d *Directorio) Nombre() string {
	return e.Nombre(d.Path)
}
