package fs

import (
	"database/sql"
	"fmt"
	u "own_wiki/system_protocol/utilidades"
	"slices"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	QUERY_LIBROS_LS       = "SELECT titulo FROM libros ORDER BY titulo"
	QUERY_CAPITULO_LS     = "SELECT CONCAT('N° ', capitulo, ') ', nombre) FROM capitulos WHERE idLibro = %d ORDER BY capitulo"
	QUERY_PAPERS_LS       = "SELECT titulo FROM papers ORDER BY titulo"
	QUERY_DISTRIBUCION_LS = "SELECT nombre FROM distribuciones WHERE tipo = '%s'"
)

const (
	QUERY_OBTENER_LIBRO = "SELECT id, titulo FROM libros WHERE titulo = '%s'"
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

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type Colecciones struct {
	Bdd    *sql.DB
	Tipo   TipoColeccion
	Indice *u.Pila[int64]
	Path   *u.Pila[string]
}

func NewColeccion(bdd *sql.DB) *Colecciones {
	return &Colecciones{
		Bdd:    bdd,
		Tipo:   TCO_COLECCION,
		Indice: u.NewPila[int64](),
		Path:   u.NewPila[string](),
	}
}

func (c *Colecciones) DeterminarRuta(ec echo.Context) error {
	path := strings.TrimSpace(ec.QueryParam("path"))

	var carpetaActual string
	var errCd error
	for subpath := range strings.SplitSeq(path, "/") {
		carpetaActual, errCd = c.Cd(subpath)
		if errCd != nil {
			break
		}
		if carpetaActual == PD_ROOT {
			return ec.Render(200, "root", DATA_ROOT)
		}
	}

	if carpetaActual == "" {
		carpetaActual = "Colecciones"
	}
	data, errLs := c.Ls(carpetaActual)

	if errCd != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Cd tuvo el error: %v", errCd), "/Colecciones"))
	}
	if errLs != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Ls tuvo el error: %v", errLs), "/Colecciones"))
	}

	return ec.Render(200, "colecciones", data)
}

func (c *Colecciones) Ls(carpetaActual string) (Data, error) {
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

		return NewData(
			NewTextoVinculo(carpetaActual, "/Root"), c.PathActual(0), opciones.Items(),
		), nil

	case TCO_LIBROS:
		query = QUERY_LIBROS_LS
	case TCO_CAPITULOS:
		if idLibro, err := c.Indice.Pick(); err == nil {
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

		return NewData(
			NewTextoVinculo(carpetaActual, "/Colecciones?path=.."), c.PathActual(0), opciones.Items(),
		), nil

	case TCO_DIST_DISCRETA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_DISCRETA)

	case TCO_DIST_CONTINUA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_CONTINUA)

	case TCO_DIST_MULTIVARIADA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_MULTIVARIADA)
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

	return NewData(
		NewTextoVinculo(carpetaActual, "/Colecciones?path=.."), c.PathActual(0), opciones.Items(),
	), nil
}

func (c *Colecciones) PathActual(profundidad int) []TextoVinculo {
	if profundidad > 2 {
		return []TextoVinculo{
			NewTextoVinculo("...", fmt.Sprintf("/Colecciones?path=%s", strings.Repeat("../", profundidad))),
		}
	}

	if elemento, err := c.Path.Desapilar(); err != nil {
		return []TextoVinculo{
			NewTextoVinculo("Own_wiki", fmt.Sprintf("/Colecciones?path=%s", strings.Repeat("../", profundidad+1))),
			NewTextoVinculo("Colecciones", fmt.Sprintf("/Colecciones?path=%s", strings.Repeat("../", profundidad))),
		}
	} else {
		textoVinculo := NewTextoVinculo(elemento, fmt.Sprintf("/Colecciones?path=%s", strings.Repeat("../", profundidad)))
		pathActual := append(c.PathActual(profundidad+1), textoVinculo)
		c.Path.Apilar(elemento)
		return pathActual
	}
}

func (c *Colecciones) Cd(subpath string) (string, error) {
	if subpath == "" {
		carpetaActual, _ := c.Path.Pick()
		return carpetaActual, nil
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
			return "", fmt.Errorf("en %s no existe la posibilidad del path dado por '%s', con posibilidades: %v", c.Tipo, subpath, posibilidades)
		} else {
			c.Tipo = TipoColeccion(posibilidades[eleccion])
			c.Path.Apilar(posibilidades[eleccion])
			return "", nil
		}
	}

	if c.Tipo != TCO_LIBROS {
		return "", fmt.Errorf("en %s no se puede buscar nada, por lo que la busqueda '%s' no tiene sentido", c.Tipo, subpath)
	}

	fila := c.Bdd.QueryRow(fmt.Sprintf(QUERY_OBTENER_LIBRO, subpath))
	var id int64
	var nombre string
	if err := fila.Scan(&id, &nombre); err != nil {
		return "", fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	c.Tipo = TCO_CAPITULOS
	c.Indice.Apilar(id)
	c.Path.Apilar(nombre)
	return subpath, nil
}

func (c *Colecciones) RutinaAtras() (string, error) {
	switch c.Tipo {
	case TCO_COLECCION:
		return PD_ROOT, nil

	case TCO_LIBROS:
		fallthrough
	case TCO_PAPERS:
		fallthrough
	case TCO_DISTRIBUCIONES:
		c.Tipo = TCO_COLECCION

	case TCO_CAPITULOS:
		_, _ = c.Indice.Desapilar()
		c.Tipo = TCO_LIBROS

	case TCO_DIST_DISCRETA:
		fallthrough
	case TCO_DIST_CONTINUA:
		fallthrough
	case TCO_DIST_MULTIVARIADA:
		c.Tipo = TCO_DISTRIBUCIONES
	}

	return c.Path.Desapilar()
}
