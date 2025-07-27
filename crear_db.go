package main

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"own_wiki/system_protocol/db"
	fs "own_wiki/system_protocol/fs"

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

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de input")
		return
	} else if len(os.Args) <= 2 {
		fmt.Println("No tiene la cantidad suficiente de argumentos, necesitas pasar el directorio de output")
		return
	}

	canalMensajes := make(chan string)
	go func(canal chan string) {
		fmt.Println("Imprimiendo mensajes")
		for {
			mensaje := <-canal
			fmt.Print(mensaje)
		}
	}(canalMensajes)

	canalInfo := make(chan db.InfoArchivos)
	canalDirectorio := make(chan fs.Directorio)
	go func(canalInfo chan db.InfoArchivos, canalMensajes chan string, dirOrigen string) {
		var infoArchivos db.InfoArchivos
		directorioRoot := fs.EstablecerDirectorio(dirOrigen, &infoArchivos, canalMensajes)

		canalMensajes <- "Procesando los archivos\n"

		var waitArchivos sync.WaitGroup
		directorioRoot.ProcesarArchivos(&waitArchivos, &infoArchivos, canalMensajes)
		waitArchivos.Wait()

		canalMensajes <- "Se termino de procesar los archivos\n"

		// Ajustar valores de info de los archivos
		infoArchivos.Incrementar()

		canalInfo <- infoArchivos
		canalDirectorio <- *directorioRoot
	}(canalInfo, canalMensajes, os.Args[1])

	canalBDD := make(chan *sql.DB)
	go func(canalBDD chan *sql.DB, canalInfo chan db.InfoArchivos) {
		bdd, err := db.EstablecerBaseDeDatos()
		if err != nil {
			fmt.Printf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
			canalBDD <- nil
			return
		}

		infoArchivos := <-canalInfo

		err = db.CrearTablas(bdd, &infoArchivos)
		if err != nil {
			fmt.Printf("No se pudo crear las tablas para la base de datos, con error: %v\n", err)
			canalBDD <- nil
			return
		}

		canalBDD <- bdd
	}(canalBDD, canalInfo)

	bdd := <-canalBDD
	if bdd == nil {
		return
	}
	defer bdd.Close()

	directorioRoot := <-canalDirectorio

	fmt.Println("Insertando datos en la base de datos")
	canalProcesamiento := make(chan func() bool)
	var bddLock sync.Mutex

	go func(bdd *sql.DB, canal chan func() bool, lock *sync.Mutex) {
		directorioRoot.InsertarDatos(bdd, &bddLock, canal)
		fmt.Println("Dejar de mandar archivos para procesar")
		close(canal)
	}(bdd, canalProcesamiento, &bddLock)

	canalEnEspera := make(chan func() bool)
	esperaNuevoProcesamiento := true
	cantidadProcesamientoPendiente := 0

	for esperaNuevoProcesamiento || cantidadProcesamientoPendiente > 0 {
		var funcProcesar func() bool
		select {
		case procesar, ok := <-canalProcesamiento:
			if !ok {
				esperaNuevoProcesamiento = false
				continue
			}
			funcProcesar = procesar

		case procesar := <-canalEnEspera:
			fmt.Println("Se intenta procesar de nuevo un archivo")
			cantidadProcesamientoPendiente--
			funcProcesar = procesar
		}

		bddLock.Lock()
		if !funcProcesar() {
			cantidadProcesamientoPendiente++
			fmt.Println("No se pudo procesar un archivo")
			canalEnEspera <- funcProcesar
		}
		bddLock.Unlock()
	}

	fmt.Println("Se termino de insertar los archivos")

	fmt.Println("Fin")
}
