package encoding

import (
	"fmt"
	"os"
	"slices"
	"sync"

	c "own_wiki/encoding/configuracion"
	"own_wiki/encoding/procesar"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	u "own_wiki/system_protocol/utilidades"

	_ "embed"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"

//go:embed tablas.json
var infoTablas string

const CANTIDAD_WORKERS = 15

var DIRECTORIOS_IGNORAR = []string{".git", ".configuracion", ".github", ".obsidian", ".trash"}

func RecorrerDirectorio(dirOrigen string, tracker *d.TrackerDependencias, canalMensajes chan string) error {
	var waitArchivos sync.WaitGroup

	canalInput := make(chan string, CANTIDAD_WORKERS)
	waitArchivos.Add(1)

	terminar := func(motivo string) {
		canalMensajes <- "Esperando a que termine por: " + motivo
		close(canalInput)
		waitArchivos.Wait()
	}

	procesarArchivo := func(path string) {
		if err := procesar.CargarArchivo(dirOrigen, path, tracker, canalMensajes); err != nil {
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
			terminar("huvo un error")
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

	terminar("termino como lo esperado")
	return nil
}

func Encodear(dirInput string, canalMensajes chan string) {
	_ = godotenv.Load()

	bddRelacional, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer b.CerrarBddRelacional(bddRelacional)

	bddNoSQL, err := b.EstablecerConexionNoSQL(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer b.CerrarBddNoSQL(bddNoSQL)
	canalMensajes <- "Se conectaron correctamente las bdd necesarias"

	tablas, err := c.CrearTablas(infoTablas)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se leyeron correctamente las tablas"

	tracker, err := d.NewTrackerDependencias(b.NewBdd(bddRelacional, bddNoSQL), tablas, canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear el tracker, se tuvo el error: %v", err)
		return
	}

	if err = RecorrerDirectorio(dirInput, tracker, canalMensajes); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo recorrer todos los archivos, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se termino el proceso de insertar datos"

	if err = tracker.TerminarProcesoInsertarDatos(); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo terminar el proceso de insertar datos, se tuvo el error: %v", err)
	} else {
		canalMensajes <- "Se termino de cargar a la base de datos"
	}
	canalMensajes <- "Se hizo la limpieza de los datos auxiliares"
}
