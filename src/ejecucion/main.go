package main

import (
	"fmt"
	"os"
	b "own_wiki/system_protocol/base_de_datos"
	c "own_wiki/system_protocol/configuracion"
	v "own_wiki/system_protocol/views"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const NOMBRE_BDD = "baseDeDatos.db"

// mdp "github.com/gomarkdown/markdown/parser"
// tp "github.com/BurntSushi/toml"
// "github.com/go-sql-driver/mysql"

func ObtenerViews(dirConfiguracion string, bdd *b.Bdd) (*v.InfoViews, error) {
	if bytes, err := os.ReadFile(fmt.Sprintf("%s/%s", dirConfiguracion, "tablas.json")); err != nil {
		return nil, fmt.Errorf("error al leer el archivo de configuracion para las tablas, con error: %v", err)

	} else if descripcionTablas, err := c.DescribirTablas(string(bytes)); err != nil {
		return nil, err

	} else {
		return c.CrearInfoViews(dirConfiguracion, descripcionTablas)
	}
}

func Visualizar(carpetaOutput, carpetaConfiguracion string, canalMensajes chan string) {
	e := echo.New()
	e.Use(middleware.Logger())

	bdd, err := b.NewBdd(carpetaOutput, NOMBRE_BDD, canalMensajes)
	if err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo establecer la conexion con la base de datos, con error: %v\n", err)
		return

	}
	defer bdd.Close()

	canalMensajes <- "Se conectaron correctamente las bdd necesarias"

	if infoViews, err := ObtenerViews(carpetaConfiguracion, bdd); err != nil {
		canalMensajes <- fmt.Sprintf("No se pudo cargar las views, con error: %v", err)

	} else {
		if err = infoViews.RegistrarRenderer(e, carpetaConfiguracion); err != nil {
			canalMensajes <- fmt.Sprintf("No se pudo crear el renderer, con error: %v", err)
			return
		}

		// Ver que hacer con esto
		e.Static("/imagenes", fmt.Sprintf("%s/%s", carpetaConfiguracion, infoViews.PathImagenes))
		e.Static("/css", fmt.Sprintf("%s/%s", carpetaConfiguracion, infoViews.PathCss))

		infoViews.GenerarEndpoints(e, bdd)
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
	var carpetaOutput string

	argumentoProcesar := 1
	for argumentoProcesar+1 < len(os.Args) {
		switch os.Args[argumentoProcesar] {
		case "-c":
			argumentoProcesar++
			carpetaConfiguracion = os.Args[argumentoProcesar]
		case "-o":
			argumentoProcesar++
			carpetaOutput = os.Args[argumentoProcesar]
		default:
			canalMensajes <- fmt.Sprintf("el argumento %s no pudo ser identificado", os.Args[argumentoProcesar])
		}
		argumentoProcesar++
	}

	configuracionValida := true
	if carpetaConfiguracion == "" {
		canalMensajes <- "Necesitas pasar el directorio de configuracion (con la flag -c)"
		configuracionValida = false
	}

	if carpetaOutput == "" {
		canalMensajes <- "Necesitas pasar el directorio de output (con la flag -o)"
		configuracionValida = false
	}

	if configuracionValida {
		Visualizar(carpetaOutput, carpetaConfiguracion, canalMensajes)
	}

	close(canalMensajes)
	waitMensajes.Wait()
}
