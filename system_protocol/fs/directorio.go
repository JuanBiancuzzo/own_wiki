package fs

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"

	"own_wiki/system_protocol/db"
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

func EstablecerDirectorio(root string, infoArchivos *db.InfoArchivos) *Directorio {
	directorioRoot := NewDirectorio(root)

	colaDirectorios := ls.NewCola[*Directorio]()
	colaDirectorios.Encolar(directorioRoot)
	fmt.Printf("El directorio para trabajar va a ser: %s\n", directorioRoot.Path)

	for !colaDirectorios.Vacia() {
		directorio, err := colaDirectorios.Desencolar()
		if err != nil {
			fmt.Printf("Se tuvo un error al operar sobre la queue con el error: %v", err)
			break
		}

		archivos, err := os.ReadDir(directorio.Path)
		if err != nil {
			fmt.Printf("Se tuvo un error al leer el directorio dando el error: %v", err)
			break
		}

		for _, archivo := range archivos {
			archivoPath := fmt.Sprintf("%s/%s", directorio.Path, archivo.Name())
			archivoValido := false

			if archivo.IsDir() && !slices.Contains(DIRECTORIOS_IGNORAR, archivo.Name()) {
				nuevoDirectorio := NewDirectorio(archivoPath)
				directorio.AgregarSubdirectorio(nuevoDirectorio)
				colaDirectorios.Encolar(nuevoDirectorio)
				archivoValido = true

			} else if !archivo.IsDir() {
				directorio.AgregarArchivo(NewArchivo(archivoPath))
				archivoValido = true
			}

			if archivoValido {
				infoArchivos.MaxPath = max(infoArchivos.MaxPath, uint32(len(archivoPath)))
			}
		}
	}

	return directorioRoot
}

func (d *Directorio) AgregarSubdirectorio(directorio *Directorio) {
	d.Subdirectorios.Push(directorio)
}

func (d *Directorio) AgregarArchivo(archivo *Archivo) {
	d.Archivos.Push(archivo)
}

func (d *Directorio) ProcesarArchivos(wg *sync.WaitGroup, infoArchivos *db.InfoArchivos) {
	for _, subdirectorio := range d.Subdirectorios.Items() {
		if subdirectorio.Vacio() {
			continue
		}

		wg.Add(1)
		go func(directorio *Directorio, wg *sync.WaitGroup) {
			directorio.ProcesarArchivos(wg, infoArchivos)
			wg.Done()
		}(subdirectorio, wg)
	}

	for _, archivo := range d.Archivos.Items() {
		archivo.Interprestarse(infoArchivos)
	}
}

func (d *Directorio) InsertarDatos(db *sql.DB, dbLock *sync.Mutex, wg *sync.WaitGroup) {
	for _, subdirectorio := range d.Subdirectorios.Items() {
		if subdirectorio.Vacio() {
			continue
		}

		wg.Add(1)
		go func(directorio *Directorio, wg *sync.WaitGroup) {
			directorio.InsertarDatos(db, dbLock, wg)
			wg.Done()
		}(subdirectorio, wg)
	}

	dbLock.Lock()
	for _, archivo := range d.Archivos.Items() {
		archivo.InsertarDatos(db)
	}
	dbLock.Unlock()
}

func (d *Directorio) Vacio() bool {
	return d.Subdirectorios.Vacia() && d.Archivos.Vacia()
}

func (d *Directorio) Nombre() string {
	separacion := strings.Split(d.Path, "/")
	return separacion[len(separacion)-1]
}

func (d *Directorio) String() string {
	resultado := fmt.Sprintf("> %s\n\t", d.Nombre())

	for _, subdirectorio := range d.Subdirectorios.Items() {
		lineas := strings.Split(subdirectorio.String(), "\n")
		representacion := strings.Join(lineas, "\n\t| ")

		resultado = fmt.Sprintf("%s%s", resultado, representacion)
	}

	for _, archivo := range d.Archivos.Items() {
		resultado = fmt.Sprintf("%s| %s\n\t", resultado, archivo.Nombre())
	}

	return resultado
}
