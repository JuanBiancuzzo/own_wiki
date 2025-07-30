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
			fmt.Println(strings.TrimSpace(mensaje))
		}
	}(canalMensajes)

	canalInfo := make(chan db.InfoArchivos)
	canalDirectorio := make(chan fs.Directorio)
	go func(canalInfo chan db.InfoArchivos, canalMensajes chan string, dirOrigen string) {
		var infoArchivos db.InfoArchivos
		var waitArchivos sync.WaitGroup

		canalMensajes <- "Creando estructura\n"
		directorioRoot := fs.EstablecerDirectorio(dirOrigen, &waitArchivos, &infoArchivos, canalMensajes)
		waitArchivos.Wait()

		// Ajustar valores de info de los archivos
		infoArchivos.Incrementar()
		canalInfo <- infoArchivos

		canalMensajes <- "Procesando los archivos\n"
		directorioRoot.RelativizarPath(fmt.Sprintf("%s/", dirOrigen))

		for archivo := range directorioRoot.IterarArchivos {
			archivo.Interprestarse(canalMensajes)
		}

		canalMensajes <- "Se termino de procesar los archivos\n"
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

	fmt.Println("Insertando datos en la base de datos")
	canalProcesamiento := make(chan e.Cargable, 100)

	go func(canal chan e.Cargable) {
		directorioRoot := <-canalDirectorio
		for archivo := range directorioRoot.IterarArchivos {
			archivo.InsertarDatos(canal)
		}
		fmt.Println("Dejar de mandar archivos para procesar")
		close(canal)
	}(canalProcesamiento)

	bdd := <-canalBDD
	if bdd == nil {
		return
	}
	defer bdd.Close()
	// bdd.SetMaxOpenConns(10)

	colaProcesar := l.NewCola[e.Cargable]()
	cantidadElementos := 0
	for cargable := range canalProcesamiento {
		if !cargable.CargarDatos(bdd, canalMensajes) {
			fmt.Printf("Encolando\n")
			cantidadElementos++
			colaProcesar.Encolar(cargable)
		}
	}

	fmt.Println("Procesando archivos, faltantes: ", cantidadElementos)
	iteracion := 0
	for !colaProcesar.Vacia() && iteracion < cantidadElementos {
		cargable, err := colaProcesar.Desencolar()
		if err != nil {
			fmt.Printf("Error al desencolar el procesamiento, con error: %v\n", err)
			break
		}
		iteracion++

		if !cargable.CargarDatos(bdd, canalMensajes) {
			fmt.Println("Encolando en el loop")
			colaProcesar.Encolar(cargable)
		} else {
			iteracion = 0
			cantidadElementos--
		}
	}

	if cantidadElementos > 0 {
		fmt.Println("Hay un error, no se pudo procesar n°", cantidadElementos, " de archivos")
	}

	fmt.Println("Se termino de insertar los archivos")

	fmt.Println("Fin")
}
