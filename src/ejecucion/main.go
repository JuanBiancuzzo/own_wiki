package main

import (
	"fmt"
	"os"
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

func Visualizar(carpetaConfiguracion string, canalMensajes chan string) {
	e := echo.New()
	e.Use(middleware.Logger())

	bdd, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return

	}
	defer bdd.Close()

	if e.Renderer, err = t.NewTemplate(carpetaConfiguracion); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo crear el renderer, con error: %v", err)
		return
	}

	e.Static("/imagenes", "ejecucion/imagenes")
	e.Static("/css", "ejecucion/css")

	fs.GenerarRutasRoot(e)

	e.GET("/Facultad", fs.NewFacultad(bdd).DeterminarRuta)
	e.GET("/Colecciones", fs.NewColeccion(bdd).DeterminarRuta)
	e.GET("/Cursos", fs.NewCursos(bdd).DeterminarRuta)

	e.Logger.Fatal(e.Start(":42069"))
}

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

	var carpetaConfiguracion string

	argumentoProcesar := 1
	for argumentoProcesar+1 < len(os.Args) {
		switch os.Args[argumentoProcesar] {
		case "-c":
			argumentoProcesar++
			carpetaConfiguracion = os.Args[argumentoProcesar]
		default:
			canalMensajes <- fmt.Sprintf("el argumento %s no pudo ser identificado", os.Args[argumentoProcesar])
		}
		argumentoProcesar++
	}

	if carpetaConfiguracion != "" {
		Visualizar(carpetaConfiguracion, canalMensajes)
	} else {
		canalMensajes <- "Necesitas pasar el directorio de configuracion (con la flag -c)"
	}

	close(canalMensajes)
	waitMensajes.Wait()
}
