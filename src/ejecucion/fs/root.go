package fs

import "fmt"

type Root struct {
}

func NewRoot() *Root {
	return &Root{}
}

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
