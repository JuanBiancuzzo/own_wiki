package encoding

import (
	"fmt"

	fs "own_wiki/encoding/fs"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	t "own_wiki/system_protocol/tablas"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"

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

	tracker, err := d.NewTrackerDependencias(b.NewBdd(bddRelacional, bddNoSQL))
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear el tracker, se tuvo el error: %v", err)
		return
	}

	tablas, err := t.NewTablas(tracker, canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear las tablas, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se crearon las tablas correctamente"

	if err = tracker.IniciarProcesoInsertarDatos(b.NewInfoArchivos(), canalMensajes); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo iniciar el proceso de insertar datos, se tuvo el error: %v", err)
		return
	}
	canalMensajes <- "Se inicio el proceso de insertar datos"

	if err = fs.RecorrerDirectorio(dirInput, tablas, canalMensajes); err != nil {
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
