package encoding

import (
	"database/sql"
	"fmt"

	fs "own_wiki/encoding/fs"
	"own_wiki/system_protocol/db"
	e "own_wiki/system_protocol/estructura"
	l "own_wiki/system_protocol/listas"

	_ "github.com/go-sql-driver/mysql"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

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

func Encodear(dirInput string, canalMensajes chan string) {
	canalInfo := make(chan db.InfoArchivos)
	canalDirectorio := make(chan fs.Root)
	go ProcesarArchivos(canalInfo, canalDirectorio, dirInput, canalMensajes)

	canalBDD := make(chan *sql.DB)
	go ConstruirBaseDeDatos(canalBDD, canalInfo, canalMensajes)

	canalMensajes <- "Insertando datos en la base de datos"
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

	if cargablesListos.Lista.Largo > 0 {
		canalMensajes <- fmt.Sprint("Hubo un error, no se procesaron: ", cargablesListos.Lista.Largo, " cargables")
	}

}
