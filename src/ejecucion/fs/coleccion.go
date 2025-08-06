package fs

import (
	"database/sql"
	"fmt"
	e "own_wiki/system_protocol/datos"
	u "own_wiki/system_protocol/utilidades"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	QUERY_LIBROS_LS       = "SELECT titulo FROM libros ORDER BY titulo"
	QUERY_CAPITULO_LS     = "SELECT CONCAT('N° ', capitulo, ') ', nombre) FROM capitulos WHERE idLibro = %d ORDER BY capitulo"
	QUERY_PAPERS_LS       = "SELECT titulo FROM papers"
	QUERY_DISTRIBUCION_LS = "SELECT nombre FROM distribuciones WHERE tipo = '%s'"
)

const (
	QUERY_OBTENER_LIBRO = "SELECT id FROM libros WHERE titulo = '%s'"
)

type TipoColeccion string

const (
	TCO_COLECCION = "Coleccion"

	TCO_LIBROS    = "Libros"
	TCO_CAPITULOS = "Capitulos"

	TCO_PAPERS = "Papers"

	TCO_DISTRIBUCIONES    = "Distribuciones"
	TCO_DIST_DISCRETA     = "Discretas"
	TCO_DIST_CONTINUA     = "Continuas"
	TCO_DIST_MULTIVARIADA = "Multivarias"
)

type Colecciones struct {
	Bdd  *sql.DB
	Tipo TipoColeccion
	Path *u.Pila[int64]
}

func NewColeccion(bdd *sql.DB) *Colecciones {
	return &Colecciones{
		Bdd:  bdd,
		Tipo: TCO_COLECCION,
		Path: u.NewPila[int64](),
	}
}

func GenerarRutaColeccion(e *echo.Echo, bdd *sql.DB) {
	colecciones := NewColeccion(bdd)

	e.GET("/Colecciones", func(c echo.Context) error {
		data := colecciones.DecidirSiguientePath(c)
		return c.Render(200, "colecciones", data)
	})
}

func (c *Colecciones) DecidirSiguientePath(ec echo.Context) Data {
	path := strings.TrimSpace(ec.QueryParam("path"))
	errCd := c.Cd(path)
	data, errLs := c.Ls()
	if errCd != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Cd tuvo el error: %v", errCd), "/Colecciones"))
	}

	if errLs != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Ls tuvo el error: %v", errLs), "/Colecciones"))
	}
	return data
}

func (c *Colecciones) Ls() (Data, error) {
	var data Data
	opciones := u.NewLista[Opcion]()

	var query string

	switch c.Tipo {
	case TCO_COLECCION:
		for _, opcion := range []string{TCO_LIBROS, TCO_PAPERS, TCO_DISTRIBUCIONES} {
			opciones.Push(
				NewOpcion(opcion, fmt.Sprintf("/Colecciones?path=%s", opcion)),
			)
		}

		return NewData(NewContenidoMinimo("Colecciones", "/Root"), opciones.Items()), nil

	case TCO_LIBROS:
		query = QUERY_LIBROS_LS
	case TCO_CAPITULOS:
		if idLibro, err := c.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_CAPITULO_LS, idLibro)
		}

	case TCO_PAPERS:
		query = QUERY_PAPERS_LS

	case TCO_DISTRIBUCIONES:
		for _, opcion := range []string{TCO_DIST_DISCRETA, TCO_DIST_CONTINUA, TCO_DIST_MULTIVARIADA} {
			opciones.Push(
				NewOpcion(opcion, fmt.Sprintf("/Colecciones?path=%s", opcion)),
			)
		}

		return NewData(NewContenidoMinimo("Colecciones", "/Colecciones?path=.."), opciones.Items()), nil

	case TCO_DIST_DISCRETA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, e.DISTRIBUCION_DISCRETA)

	case TCO_DIST_CONTINUA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, e.DISTRIBUCION_CONTINUA)

	case TCO_DIST_MULTIVARIADA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, e.DISTRIBUCION_MULTIVARIADA)
	}

	if rows, err := c.Bdd.Query(query); err != nil {
		return data, fmt.Errorf("se obtuvo un error en coleccion, al hacer query, dando el error: %v", err)

	} else {
		defer rows.Close()

		for rows.Next() {
			var nombre string
			_ = rows.Scan(&nombre)

			opciones.Push(
				NewOpcion(nombre, fmt.Sprintf("/Colecciones?path=%s", nombre)),
			)
		}
	}

	return NewData(NewContenidoMinimo("Colecciones", "/Colecciones?path=.."), opciones.Items()), nil
}

func (c *Colecciones) Cd(subpath string) error {
	if subpath == "" {
		return nil
	}

	if subpath == ".." {
		return c.RutinaAtras()
	}

	if c.Tipo == TCO_COLECCION || c.Tipo == TCO_DISTRIBUCIONES {
		var posibilidades []string
		switch c.Tipo {
		case TCO_COLECCION:
			posibilidades = []string{TCO_LIBROS, TCO_PAPERS, TCO_DISTRIBUCIONES}
		case TCO_DISTRIBUCIONES:
			posibilidades = []string{TCO_DIST_DISCRETA, TCO_DIST_CONTINUA, TCO_DIST_MULTIVARIADA}
		}

		if eleccion := slices.Index(posibilidades, subpath); eleccion < 0 {
			return fmt.Errorf("en %s no existe la posibilidad del path dado por '%s', con posibilidades: %v", c.Tipo, subpath, posibilidades)
		} else {
			c.Tipo = TipoColeccion(posibilidades[eleccion])
			return nil
		}
	}

	if c.Tipo != TCO_LIBROS {
		return fmt.Errorf("en %s no se puede buscar nada, por lo que la busqueda '%s' no tiene sentido", c.Tipo, subpath)
	}

	fila := c.Bdd.QueryRow(fmt.Sprintf(QUERY_OBTENER_LIBRO, subpath))
	var id int64
	if err := fila.Scan(&id); err != nil {
		return fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	c.Tipo = TCO_CAPITULOS
	c.Path.Apilar(id)
	return nil
}

func (c *Colecciones) RutinaAtras() error {
	switch c.Tipo {
	case TCO_COLECCION:
		return fmt.Errorf("no deberia ser posible que pongan .. aca")

	case TCO_LIBROS:
		fallthrough
	case TCO_PAPERS:
		fallthrough
	case TCO_DISTRIBUCIONES:
		c.Tipo = TCO_COLECCION

	case TCO_CAPITULOS:
		_, _ = c.Path.Desapilar()
		c.Tipo = TCO_LIBROS

	case TCO_DIST_DISCRETA:
		fallthrough
	case TCO_DIST_CONTINUA:
		fallthrough
	case TCO_DIST_MULTIVARIADA:
		c.Tipo = TCO_DISTRIBUCIONES
	}

	return nil
}
