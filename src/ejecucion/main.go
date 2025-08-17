package main

import (
	"fmt"
	"os"
	t "own_wiki/ejecucion/web_view"
	b "own_wiki/system_protocol/bass_de_datos"
	c "own_wiki/system_protocol/configuracion"
	d "own_wiki/system_protocol/dependencias"
	v "own_wiki/system_protocol/views"
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

func ObtenerTablas(dirConfiguracion string) ([]d.DescripcionTabla, error) {
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dirConfiguracion, "tablas.json")); err != nil {
		return []d.DescripcionTabla{}, fmt.Errorf("error al leer el archivo de configuracion para las tablas, con error: %v", err)

	} else {
		return c.CrearTablas(string(bytes))
	}
}

func ObtenerViews(dirConfiguracion string, bdd *b.Bdd, tablas []d.DescripcionTabla) (*v.InfoViews, error) {
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dirConfiguracion, "views.json")); err != nil {
		return nil, fmt.Errorf("error al leer el archivo de configuracion para las views, con error: %v", err)

	} else {
		return c.CrearInfoViews(string(bytes))
	}
}

func Visualizar(carpetaConfiguracion string, canalMensajes chan string) {
	e := echo.New()
	e.Use(middleware.Logger())

	bddRelacional, err := b.EstablecerConexionRelacional(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return

	}
	defer bddRelacional.Close()

	bddNoSQL, err := b.EstablecerConexionNoSQL(canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return
	}
	defer b.CerrarBddNoSQL(bddNoSQL)

	bdd := b.NewBdd(bddRelacional, bddNoSQL)
	canalMensajes <- "Se conectaron correctamente las bdd necesarias"

	if tablas, err := ObtenerTablas(carpetaConfiguracion); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo cargar las tablas, con error: %v", err)

	} else if infoViews, err := ObtenerViews(carpetaConfiguracion, bdd, tablas); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo cargar las views, con error: %v", err)

	} else {
		carpetaTemplates := fmt.Sprintf("%s/%s", carpetaConfiguracion, infoViews.PathTemplates)

		if e.Renderer, err = t.NewTemplate(carpetaTemplates); err != nil {
			canalMensajes <- fmt.Sprintf("No se pudo crear el renderer, con error: %v", err)
			return
		}

		// Ver que hacer con esto
		e.Static("/imagenes", "ejecucion/imagenes")
		e.Static("/css", fmt.Sprintf("%s/%s", carpetaConfiguracion, infoViews.PathCss))

		infoViews.GenerarEndpoints(e)
		e.Logger.Fatal(e.Start(":42069"))
	}
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
