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
	Subdirectorios map[string]*Directorio
	Archivos       *ls.Lista[*Archivo]
}

func NewRoot(path string) *Directorio {
	return &Directorio{
		Padre:          nil,
		Path:           path,
		Subdirectorios: make(map[string]*Directorio),
		Archivos:       ls.NewLista[*Archivo](),
	}
}

func NewDirectorio(padre *Directorio, path string) *Directorio {
	return &Directorio{
		Padre:          padre,
		Path:           path,
		Subdirectorios: make(map[string]*Directorio),
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

		wg.Add(1)
		go func(lista *ls.Lista[string], directorio *Directorio, infoArchivo *db.InfoArchivos) {
			for archivoPath := range lista.Iterar {
				if archivo, err := NewArchivo(directorio, archivoPath, infoArchivos, canal); err != nil {
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

/*
	path := "hola/tanto/tiempo/archivo.md"

	carpetas :=
		root
		  |-> hola
			  |-> tanto
				  |-> tiempo
					  |-> archivo.md
		  |-> chau
			  |-> chua2

	(*chua2).EncontrarArchivo(path)
*/

func (d *Directorio) EncontrarArchivo(path string) (*Archivo, error) {
	if strings.Contains(path, d.Path) {
		pathRelativo := path[len(d.Path)+1:]
		indiceSlash := strings.Index(pathRelativo, "/")
		if indiceSlash >= 0 {
			if subdirectorio, ok := d.Subdirectorios[pathRelativo[:indiceSlash]]; !ok {
				return nil, fmt.Errorf("no existe la carpeta dada por %s", pathRelativo[:indiceSlash])
			} else {
				return subdirectorio.EncontrarArchivo(path)
			}
		}

		for archivo := range d.Archivos.Iterar {
			if archivo.Path == path {
				return archivo, nil
			}
		}

	} else if d.Padre != nil {
		return d.Padre.EncontrarArchivo(path)
	}

	return nil, fmt.Errorf("no existe el archivo con ese path")
}

func (d *Directorio) RelativizarPath(path string) {
	if d.Padre == nil {
		d.Path = "/"
	} else {
		d.Path = strings.Replace(d.Path, path, "", 1)
	}

	for archivo := range d.Archivos.Iterar {
		archivo.RelativizarPath(path)
	}

	for _, subdirectorio := range d.Subdirectorios {
		subdirectorio.RelativizarPath(path)
	}
}

func (d *Directorio) AgregarSubdirectorio(nombreDirectorio string, directorio *Directorio) {
	d.Subdirectorios[nombreDirectorio] = directorio
}

func (d *Directorio) AgregarArchivo(archivo *Archivo) {
	d.Archivos.Push(archivo)
}

func (d *Directorio) IterarArchivos(yield func(*Archivo) bool) {
	directorios := ls.NewCola[*Directorio]()
	directorios.Encolar(d)

	for !directorios.Vacia() {
		directorio, err := directorios.Desencolar()
		if err != nil {
			return
		}

		for archivo := range directorio.Archivos.Iterar {
			if !yield(archivo) {
				return
			}
		}

		for _, subdirectorio := range directorio.Subdirectorios {
			directorios.Encolar(subdirectorio)
		}
	}
}

func (d *Directorio) Nombre() string {
	return e.Nombre(d.Path)
}
