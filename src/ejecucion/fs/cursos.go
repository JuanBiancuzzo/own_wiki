package fs

import (
	"database/sql"
	"fmt"
)

type Cursos struct {
	Bdd *sql.DB
}

func NewCursos(bdd *sql.DB) *Cursos {
	return &Cursos{
		Bdd: bdd,
	}
}

func (c *Cursos) Ls() ([]string, error) {
	return []string{}, nil
}

func (c *Cursos) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
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
