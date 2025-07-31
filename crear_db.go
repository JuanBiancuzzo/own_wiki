package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
	fs "own_wiki/system_protocol/fs"
	l "own_wiki/system_protocol/listas"

	_ "github.com/go-sql-driver/mysql"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

/*
Idea actual para probar
Vamos a recolectar todos los archivos de mi obsidian, y si quiero crear un archivo, lo hago por medio de este
script, de esa forma voy creando todo lo que necesite desde aca y mientras uso obsidian para renderizarlo

Ideas:
Referencias
Estructura de arbol -> Esto hace las colecciones, las facultad, las investigaciones, los cursos y los proyectos
  - Seccion
  - Notas
  - Es una lista de posibles archivos
  - Esto tambien puede ser una seccion para crear
*/

func ImprimirMensajes(canal chan string, wg *sync.WaitGroup) {
	for mensaje := range canal {
		fmt.Println(strings.TrimSpace(mensaje))
	}
	wg.Done()
}

func ProcesarArchivos(canalInfo chan db.InfoArchivos, canalDirectorio chan fs.Root, dirOrigen string, canalMensajes chan string) {
	var infoArchivos db.InfoArchivos

	canalMensajes <- "Creando estructura\n"
	directorioRoot := fs.EstablecerDirectorio(dirOrigen, &infoArchivos, canalMensajes)

	// Ajustar valores de info de los archivos
	infoArchivos.Incrementar()
	canalInfo <- infoArchivos

	canalMensajes <- "Se termino de procesar los archivos\n"
	canalDirectorio <- *directorioRoot
}

func ConstruirBaseDeDatos(canalBDD chan *sql.DB, canalInfo chan db.InfoArchivos, canalMensajes chan string) {
	bdd, err := db.EstablecerBaseDeDatos()
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}

	infoArchivos := <-canalInfo

	err = db.CrearTablas(bdd, &infoArchivos)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas para la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}

	canalBDD <- bdd
}

func EvaluarCargable(bdd *sql.DB, canalMensajes chan string, cargable e.Cargable, cola *l.Cola[e.Cargable]) {
	if id, err := cargable.CargarDatos(bdd, canalMensajes); err == nil {
		for _, cargable := range cargable.ResolverDependencias(id) {
			cola.Encolar(cargable)
		}

	} else {
		canalMensajes <- fmt.Sprintf("Error al cargar: %v", err)
	}
}

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de input")
		return
	} else if len(os.Args) <= 2 {
		fmt.Println("No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de output")
		return
	}

	canalMensajes := make(chan string, 100)
	var waitMensajes sync.WaitGroup
	waitMensajes.Add(1)
	go ImprimirMensajes(canalMensajes, &waitMensajes)

	canalInfo := make(chan db.InfoArchivos)
	canalDirectorio := make(chan fs.Root)
	go ProcesarArchivos(canalInfo, canalDirectorio, os.Args[1], canalMensajes)

	canalBDD := make(chan *sql.DB)
	go ConstruirBaseDeDatos(canalBDD, canalInfo, canalMensajes)

	fmt.Println("Insertando datos en la base de datos")
	canalIndependientes := make(chan e.Cargable, 100)

	go func(canal chan e.Cargable, canalMensajes chan string) {
		root := <-canalDirectorio
		for _, archivo := range root.Archivos {
			archivo.EstablecerDependencias(canal, canalMensajes)
		}
		canalMensajes <- "Dejar de mandar archivos para procesar"
		close(canal)
	}(canalIndependientes, canalMensajes)

	bdd := <-canalBDD
	if bdd == nil {
		return
	}
	defer bdd.Close()
	// bdd.SetMaxOpenConns(10)

	canalMensajes <- "Cargando los archivos sin dependencias"

	cargablesListos := l.NewCola[e.Cargable]()
	for cargable := range canalIndependientes {
		EvaluarCargable(bdd, canalMensajes, cargable, cargablesListos)
	}

	canalMensajes <- "Cargados todos los archivos sin dependencias, ahora procesando los que tengan dependencias"

	for !cargablesListos.Vacia() {
		if cargable, err := cargablesListos.Desencolar(); err != nil {
			canalMensajes <- fmt.Sprintf("Error al desencolar el procesamiento, con error: %v", err)
			break

		} else {
			EvaluarCargable(bdd, canalMensajes, cargable, cargablesListos)
		}
	}

	close(canalMensajes)
	waitMensajes.Wait()

	if cargablesListos.Lista.Largo > 0 {
		fmt.Println("Hubo un error, no se procesaron: ", cargablesListos.Lista.Largo, " cargables")
	}

	fmt.Println("Se termino de insertar los archivos")
	fmt.Println("Fin")
}
