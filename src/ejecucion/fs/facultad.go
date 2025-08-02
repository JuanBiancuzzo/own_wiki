package fs

import (
	"database/sql"
	"fmt"

	l "own_wiki/system_protocol/listas"
)

const QUERY_CARRERAS_LS = "SELECT nombre FROM carreras"

const (
	QUERY_EXISTENCIA_CARRERA = "SELECT id, nombre FROM carreras WHERE nombre = ?"
	QUERY_EXISTENCIA_MATERIA = "SELECT res.id, res.nombre FROM (SELECT id, nombre FROM materias INNER JOIN carreras ON materias.idCarrera = carreras.id) AS res WHERE res.nombre = ?"
	QUERY_EXISTENCIA_TEMA    = "SELECT res.id, res.nombre FROM (SELECT id, nombre FROM temasMateria INNER JOIN materias ON temasMateria.idMateria = materias.id) AS res WHERE res.nombre = ?"
)

type TipoFacultad byte

const (
	TF_CARRERA = iota
	TF_MATERIA
	TF_TEMA
)

type IdFacultad struct {
	Id     int64
	Nombre string
}

type Facultad struct {
	Bdd  *sql.DB
	Tipo TipoFacultad
	Path *l.Pila[IdFacultad]
}

func NewFacultad(bdd *sql.DB) *Facultad {
	return &Facultad{
		Bdd:  bdd,
		Tipo: TF_CARRERA,
		Path: l.NewPila[IdFacultad](),
	}
}

func (f *Facultad) Ls() ([]string, error) {
	switch f.Tipo {
	case TF_CARRERA:
		if rows, err := f.Bdd.Query(QUERY_CARRERAS_LS); err != nil {
			return []string{}, fmt.Errorf("se obtuvo un error en facultad, al hacer query de carrera, dando el error: %v", err)

		} else {
			columnas := l.NewLista[string]()
			defer rows.Close()
			for rows.Next() {
				var nombreMateria string
				_ = rows.Scan(&nombreMateria)
				columnas.Push(nombreMateria)
			}

			return columnas.Items(), nil
		}
	}

	return []string{}, nil
}

func (f *Facultad) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		_, _ = f.Path.Desapilar()
		switch f.Tipo {
		case TF_CARRERA:
			return cache.ObtenerSubpath(PD_ROOT)

		case TF_MATERIA:
			f.Tipo = TF_CARRERA
		case TF_TEMA:
			f.Tipo = TF_MATERIA
		}

		return f, nil
	}

	var query string
	switch f.Tipo {
	case TF_CARRERA:
		query = QUERY_EXISTENCIA_CARRERA
	case TF_MATERIA:
		query = QUERY_EXISTENCIA_MATERIA
	case TF_TEMA:
		query = QUERY_EXISTENCIA_TEMA
	}

	fila := f.Bdd.QueryRow(query, subpath)
	var idElemento int64
	if err := fila.Scan(&idElemento); err != nil {
		return f, fmt.Errorf("no existe posible solucion para el cd a '%s'", subpath)
	}

	f.Path.Apilar(IdFacultad{Id: idElemento, Nombre: subpath})
	return f, nil
}
