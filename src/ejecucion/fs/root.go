package fs

import (
	"github.com/labstack/echo/v4"
)

const (
	PD_ROOT      = "Root"
	PD_FACULTAD  = "Facultad"
	PD_CURSOS    = "Cursos"
	PD_COLECCION = "Colecciones"
)

var DATA_ROOT = NewData(NewTextoVinculo("Elecciones", "/Root"),
	[]TextoVinculo{},
	[]Opcion{
		NewOpcion(PD_FACULTAD, "/Facultad"),
		NewOpcion(PD_CURSOS, "/Cursos"),
		NewOpcion(PD_COLECCION, "/Colecciones"),
	},
)

func GenerarRutasRoot(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", DATA_ROOT)
	})

	e.GET("/Root", func(c echo.Context) error {
		return c.Render(200, "root", DATA_ROOT)
	})
}
