package fs

import (
	"github.com/labstack/echo/v4"
)

const (
	PD_FACULTAD  = "Facultad"
	PD_CURSOS    = "Cursos"
	PD_COLECCION = "Colecciones"
)

func GenerarRutasRoot(e *echo.Echo) {
	data := NewData(NewContenidoMinimo("Elecciones", "/Root"),
		[]Opcion{
			NewOpcion(PD_FACULTAD, "/Facultad"),
			NewOpcion(PD_CURSOS, "/Cursos"),
			NewOpcion(PD_COLECCION, "/Colecciones"),
		},
	)

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", data)
	})

	e.GET("/Root", func(c echo.Context) error {
		return c.Render(200, "root", data)
	})
}
