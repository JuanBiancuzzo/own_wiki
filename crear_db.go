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

	canalDB := make(chan *sql.DB)
	go func(canalDB chan *sql.DB, canalInfo chan db.InfoArchivos) {
		baseDeDatos, err := db.EstablecerBaseDeDatos()
		if err != nil {
			fmt.Printf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
			canalDB <- nil
			return
		}

		infoArchivos := <-canalInfo

		err = db.CrearTablas(baseDeDatos, &infoArchivos)
		if err != nil {
			fmt.Printf("No se pudo crear las tablas para la base de datos, con error: %v\n", err)
			canalDB <- nil
			return
		}

		canalDB <- baseDeDatos
	}(canalDB, canalInfo)

	baseDeDatos := <-canalDB
	if baseDeDatos == nil {
		return
	}
	defer baseDeDatos.Close()

	directorioRoot := <-canalDirectorio

	fmt.Println("Insertando datos en la base de datos")
	var waitInsersion sync.WaitGroup
	var dbLock sync.Mutex

	directorioRoot.InsertarDatos(baseDeDatos, &dbLock, &waitInsersion)

	waitInsersion.Wait()
	fmt.Println("Se termino de insertar los archivos")

	fmt.Println("Fin")
}
