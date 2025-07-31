package fs

import (
	"strings"

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
