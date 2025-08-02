package fs

import (
	"database/sql"
	"fmt"
)

type Colecciones struct {
	Bdd *sql.DB
}

func NewColeccion(bdd *sql.DB) *Colecciones {
	return &Colecciones{
		Bdd: bdd,
	}
}

func (c *Colecciones) Ls() ([]string, error) {
	return []string{}, nil
}

func (c *Colecciones) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		return cache.ObtenerSubpath(PD_ROOT)
	}

	switch subpath {
	case PD_FACULTAD:
		fallthrough
	case PD_CURSOS:
		fallthrough
	case PD_COLECCION:
	}

	return c, fmt.Errorf("no existe posible solucion para el cd a '%s'", subpath)
}
