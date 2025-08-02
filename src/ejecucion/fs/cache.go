package fs

import (
	"database/sql"
	"fmt"
)

type EstadoDirectorio string

const (
	PD_ROOT      = "/"
	PD_FACULTAD  = "Facultad"
	PD_CURSOS    = "Cursos"
	PD_COLECCION = "Colecciones"
)

type Cache struct {
	Bdd      *sql.DB
	Subpaths map[EstadoDirectorio]Subpath
}

func NewCache(bdd *sql.DB) *Cache {
	return &Cache{
		Bdd:      bdd,
		Subpaths: make(map[EstadoDirectorio]Subpath),
	}
}

func (c *Cache) ObtenerSubpath(estado EstadoDirectorio) (Subpath, error) {
	if subpath, ok := c.Subpaths[estado]; ok {
		return subpath, nil
	}

	var nuevoEstado Subpath
	switch estado {
	case PD_ROOT:
		nuevoEstado = NewRoot()
	case PD_FACULTAD:
		nuevoEstado = NewFacultad(c.Bdd)
	case PD_CURSOS:
		nuevoEstado = NewCursos(c.Bdd)
	case PD_COLECCION:
		nuevoEstado = NewColeccion(c.Bdd)
	default:
		return nuevoEstado, fmt.Errorf("de alguna forma se paso el estado '%s' que no esta registrado", estado)
	}

	return nuevoEstado, nil
}
