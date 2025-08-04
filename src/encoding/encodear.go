package encoding

import (
	"database/sql"
	"fmt"
	"sync"
	"unsafe"

	fs "own_wiki/encoding/fs"
	b "own_wiki/system_protocol/baseDeDatos"
	e "own_wiki/system_protocol/datos"
	l "own_wiki/system_protocol/utilidades"

	ts "github.com/tree-sitter/go-tree-sitter"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func ProcesarArchivos(canalInfo chan b.InfoArchivos, canalDirectorio chan fs.Root, dirOrigen string, canalMensajes chan string) {
	var infoArchivos b.InfoArchivos

	canalMensajes <- "Creando estructura\n"
	directorioRoot := fs.EstablecerDirectorio(dirOrigen, &infoArchivos, canalMensajes)

	// Ajustar valores de info de los archivos
	infoArchivos.Incrementar()
	canalInfo <- infoArchivos
	close(canalInfo)

	canalMensajes <- "Se termino de procesar los archivos\n"
	canalDirectorio <- *directorioRoot
	close(canalDirectorio)
}

func ConstruirBaseRelacional(canalBDD chan *sql.DB, canalInfo chan b.InfoArchivos, canalMensajes chan string) {
	bdd, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}
	infoArchivos := <-canalInfo

	err = b.CrearTablas(bdd, &infoArchivos)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas para la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}

	canalBDD <- bdd
	close(canalBDD)
}

func ConstruirBaseNoSQL(canalBDD chan *mongo.Database, canalMensajes chan string) {
	bdd, err := b.EstablecerConexionNoSQL(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}

	err = b.CrearColecciones(bdd)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las colecciones para la base de datos, con error: %v\n", err)
		canalBDD <- nil
		return
	}

	canalBDD <- bdd
	close(canalBDD)
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

func CargarDatos(bddRelacional *sql.DB, canalIndependiente chan e.Cargable, wg *sync.WaitGroup, canalMensajes chan string) {
	canalMensajes <- "Cargando los archivos sin dependencias"

	cargablesListos := l.NewCola[e.Cargable]()
	for cargable := range canalIndependiente {
		EvaluarCargable(bddRelacional, canalMensajes, cargable, cargablesListos)
	}

	canalMensajes <- "Cargados todos los archivos sin dependencias, ahora procesando los que tengan dependencias"

	for cargable := range cargablesListos.DesencolarIterativamente {
		EvaluarCargable(bddRelacional, canalMensajes, cargable, cargablesListos)
	}

	if cargablesListos.Lista.Largo > 0 {
		canalMensajes <- fmt.Sprint("Hubo un error, no se procesaron: ", cargablesListos.Lista.Largo, " cargables")
	}

	wg.Done()
}

func CargarDocumentos(bddNoSQL *mongo.Database, canalIndependiente chan e.A, wg *sync.WaitGroup, canalMensajes chan string) {
	canalMensajes <- "Cargando los documentos"
	wg.Done()
}

func Encodear(dirInput string, canalMensajes chan string) {

	code := []byte("const foo = 1 + 2")

	parser := ts.NewParser()
	defer parser.Close()

	if pointer, err := tsb.LanguageJavascript(); err != nil {
		canalMensajes <- fmt.Sprintf("Error al cargar libreria, con el error: %v", err)
	} else {
		parser.SetLanguage(ts.NewLanguage(unsafe.Pointer(pointer)))

		tree := parser.Parse(code, nil)
		defer tree.Close()

		root := tree.RootNode()
		fmt.Println(root.ToSexp())
	}

	canalInfo := make(chan b.InfoArchivos)
	canalDirectorio := make(chan fs.Root)
	go ProcesarArchivos(canalInfo, canalDirectorio, dirInput, canalMensajes)

	_ = godotenv.Load()

	canalBddRelacional := make(chan *sql.DB)
	go ConstruirBaseRelacional(canalBddRelacional, canalInfo, canalMensajes)

	canalBddNoSQL := make(chan *mongo.Database)
	go ConstruirBaseNoSQL(canalBddNoSQL, canalMensajes)

	canalDatos := make(chan e.Cargable, 100)
	canalDocumentos := make(chan e.A, 100)

	go func(canalDatos chan e.Cargable, canalDocumentos chan e.A, canalMensajes chan string) {
		root := <-canalDirectorio
		for _, archivo := range root.Archivos {
			archivo.EstablecerDependencias(canalDatos, canalDocumentos, canalMensajes)
		}

		canalMensajes <- "Dejar de mandar archivos para procesar"
		close(canalDatos)
		close(canalDocumentos)

	}(canalDatos, canalDocumentos, canalMensajes)

	bddRelacional := <-canalBddRelacional
	defer b.CerrarBddRelacional(bddRelacional)

	bddNoSQL := <-canalBddNoSQL
	defer b.CerrarBddNoSQL(bddNoSQL)

	if bddRelacional == nil || bddNoSQL == nil {
		canalMensajes <- "No se pudo hacer una conexion con las bases de datos"
		return
	}
	canalMensajes <- "Insertando datos en la base de datos"

	var waitCarga sync.WaitGroup

	waitCarga.Add(1)
	go CargarDatos(bddRelacional, canalDatos, &waitCarga, canalMensajes)

	waitCarga.Add(1)
	go CargarDocumentos(bddNoSQL, canalDocumentos, &waitCarga, canalMensajes)

	waitCarga.Wait()
	canalMensajes <- "Se termino de cargar a la base de datos"
}
