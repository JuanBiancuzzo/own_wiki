package main

import (
	"fmt"
	"own_wiki/ejecucion/fs"
	t "own_wiki/ejecucion/web_view"
	b "own_wiki/system_protocol/bass_de_datos"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func main() {
	_ = godotenv.Load()
	var waitMensajes sync.WaitGroup
	canalMensajes := make(chan string, 100)

	waitMensajes.Add(1)
	go func(canal chan string, wg *sync.WaitGroup) {
		for mensaje := range canal {
			fmt.Println(strings.TrimSpace(mensaje))
		}
		wg.Done()
	}(canalMensajes, &waitMensajes)

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

	e.GET("/Facultad", fs.NewFacultad(bdd).DeterminarRuta)
	e.GET("/Colecciones", fs.NewColeccion(bdd).DeterminarRuta)
	e.GET("/Cursos", fs.NewCursos(bdd).DeterminarRuta)

	e.Logger.Fatal(e.Start(":42069"))

	close(canalMensajes)
	waitMensajes.Wait()
}
