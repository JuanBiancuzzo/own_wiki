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
	Padre          *Directorio
	Path           string
	Subdirectorios *ls.Lista[*Directorio]
	Archivos       *ls.Lista[*Archivo]
}

func NewRoot(path string) *Directorio {
	return &Directorio{
		Padre:          nil,
		Path:           path,
		Subdirectorios: ls.NewLista[*Directorio](),
		Archivos:       ls.NewLista[*Archivo](),
	}
}

func NewDirectorio(padre *Directorio, path string) *Directorio {
	return &Directorio{
		Padre:          padre,
		Path:           path,
		Subdirectorios: ls.NewLista[*Directorio](),
		Archivos:       ls.NewLista[*Archivo](),
	}
}

func EstablecerDirectorio(root string, wg *sync.WaitGroup, infoArchivos *db.InfoArchivos, canal chan string) *Directorio {
	directorioRoot := NewRoot(root)

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
			archivoPath := fmt.Sprintf("%s/%s", directorio.Path, archivo.Name())

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, archivo.Name()) {
				subDirectorio := NewDirectorio(directorio, archivoPath)
				directorio.AgregarSubdirectorio(subDirectorio)
				colaDirectorios.Encolar(subDirectorio)

			} else if !archivo.IsDir() {
				listaArchivos.Push(archivoPath)
			}
		}

		wg.Add(1)
		go func(lista *ls.Lista[string], directorio *Directorio, infoArchivo *db.InfoArchivos) {
			for archivoPath := range lista.Iterar {
				if archivo, err := NewArchivo(directorio, archivoPath, infoArchivos); err != nil {
					canal <- fmt.Sprintf("Se tuvo un error al crear un archivo, con error: %v", err)

				} else {
					directorio.AgregarArchivo(archivo)
				}
			}
			wg.Done()
		}(listaArchivos, directorio, infoArchivos)
	}

	return directorioRoot
}

func (d *Directorio) ArchivoMasCercanoConTag(tag string) (*Archivo, error) {
	for archivo := range d.Archivos.Iterar {
		if slices.Contains(archivo.Meta.Tags, tag) {
			return archivo, nil
		}
	}

	if d.Padre == nil {
		return nil, fmt.Errorf("se llegó al root y no se encontró archivo con tag: %s", tag)
	}
	return d.Padre.ArchivoMasCercanoConTag(tag)
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
