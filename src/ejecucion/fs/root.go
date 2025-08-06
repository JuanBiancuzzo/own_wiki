package fs

import (
	"fmt"

	"github.com/labstack/echo/v4"
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

func GenerarRutaResto(e *echo.Echo) {
	e.GET("/Facultad", func(c echo.Context) error {
		opcion := NewContenidoMinimo("Facultad", "/Root")
		return c.Render(200, "otros", opcion)
	})

	e.GET("/Cursos", func(c echo.Context) error {
		opcion := NewContenidoMinimo("Cursos", "/Root")
		return c.Render(200, "otros", opcion)
	})

}

type Root struct {
	Echo *echo.Echo
}

func NewRoot(echo *echo.Echo) *Root {
	return &Root{
		Echo: echo,
	}
}

func (r *Root) Algo() {}

func (r *Root) Ls() ([]string, error) {
	return []string{PD_FACULTAD, PD_CURSOS, PD_COLECCION}, nil
}

func (r *Root) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		return r, nil
	}

	switch subpath {
	case PD_FACULTAD:
		fallthrough
	case PD_CURSOS:
		fallthrough
	case PD_COLECCION:
		return cache.ObtenerSubpath(EstadoDirectorio(subpath))
	}

	return r, fmt.Errorf("no existe posible solucion para el cd a '%s'", subpath)
}
