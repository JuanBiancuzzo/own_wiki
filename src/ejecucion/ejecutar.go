package ejecucion

import (
	"fmt"
	"own_wiki/ejecucion/fs"
	t "own_wiki/ejecucion/web_view"
	b "own_wiki/system_protocol/bass_de_datos"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func Ejecutar(canalMensajes chan string) {
	_ = godotenv.Load()

	e := echo.New()
	e.Use(middleware.Logger())

	bdd, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return

	}
	defer bdd.Close()

	e.Renderer = t.NewTemplate()
	e.Static("/imagenes", "ejecucion/imagenes")
	e.Static("/css", "ejecucion/css")

	fs.GenerarRutasRoot(e)
	fs.GenerarRutaColeccion(e, bdd)
	fs.GenerarRutaFacultad(e, bdd)
	fs.GenerarRutaCursos(e, bdd)

	e.Logger.Fatal(e.Start(":42069"))
}
