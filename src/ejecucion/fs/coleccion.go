package fs

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/utilidades"
	"slices"
)

const (
	QUERY_LIBROS_LS       = "SELECT titulo FROM libros"
	QUERY_CAPITULO_LS     = "SELECT CONCAT('N° ', capitulo, ') ', nombre) FROM capitulos WHERE idLibro = %d"
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

	TCO_DISTRIBUCIONES  = "Distribuciones"
	TCO_DIST_DISCRETA   = "Discretas"
	TCO_DIST_CONTINUA   = "Continuas"
	TCO_DIST_MULTIVARIA = "Multivarias"
)

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type Colecciones struct {
	Bdd  *sql.DB
	Tipo TipoColeccion
	Path *l.Pila[int64]
}

func NewColeccion(bdd *sql.DB) *Colecciones {
	return &Colecciones{
		Bdd:  bdd,
		Tipo: TCO_COLECCION,
		Path: l.NewPila[int64](),
	}
}

func (c *Colecciones) Ls() ([]string, error) {
	var query string

	switch c.Tipo {
	case TCO_COLECCION:
		return []string{TCO_LIBROS, TCO_PAPERS, TCO_DISTRIBUCIONES}, nil

	case TCO_LIBROS:
		query = QUERY_LIBROS_LS
	case TCO_CAPITULOS:
		if idLibro, err := c.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_CAPITULO_LS, idLibro)
		}

	case TCO_PAPERS:
		query = QUERY_PAPERS_LS

	case TCO_DISTRIBUCIONES:
		return []string{TCO_DIST_DISCRETA, TCO_DIST_CONTINUA, TCO_DIST_MULTIVARIA}, nil

	case TCO_DIST_DISCRETA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_DISCRETA)

	case TCO_DIST_CONTINUA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_CONTINUA)

	case TCO_DIST_MULTIVARIA:
		query = fmt.Sprintf(QUERY_DISTRIBUCION_LS, DISTRIBUCION_MULTIVARIADA)
	}

	if rows, err := c.Bdd.Query(query); err != nil {
		return []string{}, fmt.Errorf("se obtuvo un error en coleccion, al hacer query, dando el error: %v", err)

	} else {
		columnas := l.NewLista[string]()
		defer rows.Close()
		for rows.Next() {
			var nombre string
			_ = rows.Scan(&nombre)
			columnas.Push(nombre)
		}

		return columnas.Items(), nil
	}
}

func (c *Colecciones) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		return c.RutinaAtras(cache)
	}

	if c.Tipo == TCO_COLECCION || c.Tipo == TCO_DISTRIBUCIONES {
		var posibilidades []string
		switch c.Tipo {
		case TCO_COLECCION:
			posibilidades = []string{TCO_LIBROS, TCO_PAPERS, TCO_DISTRIBUCIONES}
		case TCO_DISTRIBUCIONES:
			posibilidades = []string{TCO_DIST_DISCRETA, TCO_DIST_CONTINUA, TCO_DIST_MULTIVARIA}
		}

		if eleccion := slices.Index(posibilidades, subpath); eleccion < 0 {
			return c, fmt.Errorf("en %s no existe la posibilidad del path dado por '%s'", c.Tipo, subpath)
		} else {
			c.Tipo = TipoColeccion(posibilidades[eleccion])
			return c, nil
		}
	}

	if c.Tipo != TCO_LIBROS {
		return c, fmt.Errorf("en %s no se puede buscar nada, por lo que la busqueda '%s' no tiene sentido", c.Tipo, subpath)
	}

	fila := c.Bdd.QueryRow(fmt.Sprintf(QUERY_OBTENER_LIBRO, subpath))
	var id int64
	if err := fila.Scan(&id); err != nil {
		return c, fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	c.Tipo = TCO_CAPITULOS
	c.Path.Apilar(id)
	return c, nil
}

func (c *Colecciones) RutinaAtras(cache *Cache) (Subpath, error) {
	switch c.Tipo {
	case TCO_COLECCION:
		return cache.ObtenerSubpath(PD_ROOT)

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
	case TCO_DIST_MULTIVARIA:
		c.Tipo = TCO_DISTRIBUCIONES
	}

	return c, nil
}
