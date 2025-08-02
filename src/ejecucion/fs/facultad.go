package fs

import (
	"database/sql"
	"fmt"

	l "own_wiki/system_protocol/listas"
)

const (
	QUERY_CARRERAS_LS = "SELECT nombre FROM carreras"
	QUERY_MATERIAS_LS = "SELECT materias.nombre FROM materias INNER JOIN (SELECt id FROM carreras WHERE id = %d) AS carrera ON materias.idCarrera = carrera.id"
	QUERY_TEMAS_LS    = "SELECT temasMateria.nombre FROM temasMateria INNER JOIN (SELECt id FROM materias WHERE id = %d) AS materia ON temasMateria.idMateria = materia.id"
)

const (
	QUERY_EXISTENCIA_CARRERA = "SELECT id, nombre FROM carreras WHERE nombre = '%s'"
	QUERY_EXISTENCIA_MATERIA = "SELECT materias.id, materias.nombre FROM materias INNER JOIN (SELECt id FROM carreras WHERE id = %d) AS carrera ON materias.idCarrera = carrera.id WHERE materias.nombre = '%s'"
	QUERY_EXISTENCIA_TEMA    = "SELECT temasMateria.id, temasMateria.nombre FROM temasMateria INNER JOIN (SELECt id FROM materias WHERE id = %d) AS materia ON temasMateria.idMateria = materia.id WHERE temasMateria.nombre = '%s'"
)

type TipoFacultad byte

const (
	TF_CARRERA = iota
	TF_MATERIA
	TF_TEMA
	TF_NOTA
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
	var query string
	switch f.Tipo {
	case TF_CARRERA:
		query = QUERY_CARRERAS_LS
	case TF_MATERIA:
		if idCarrera, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_MATERIAS_LS, idCarrera.Id)
		}
	case TF_TEMA:
		if idMateria, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_TEMAS_LS, idMateria.Id)
		}
	}

	if rows, err := f.Bdd.Query(query); err != nil {
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

	var query string = ""
	switch f.Tipo {
	case TF_CARRERA:
		query = fmt.Sprintf(QUERY_EXISTENCIA_CARRERA, subpath)
	case TF_MATERIA:
		if idCarrera, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_EXISTENCIA_MATERIA, idCarrera.Id, subpath)
		}
	case TF_TEMA:
		if idMateria, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_EXISTENCIA_TEMA, idMateria.Id, subpath)
		}
	case TF_NOTA:
		return f, fmt.Errorf("ya se esta viendo todos los archivos, no hay subcarpetas")
	}

	fmt.Printf("QUERY: %s\n", query)
	if query == "" {
		return f, fmt.Errorf("hubo un error en la query, y esta vacia")
	}

	fila := f.Bdd.QueryRow(query)
	var id int64
	var nombre string

	if err := fila.Scan(&id, &nombre); err != nil {
		return f, fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	switch f.Tipo {
	case TF_CARRERA:
		f.Tipo = TF_MATERIA
	case TF_MATERIA:
		f.Tipo = TF_TEMA
	case TF_TEMA:
		f.Tipo = TF_NOTA
	}

	f.Path.Apilar(IdFacultad{Id: id, Nombre: nombre})
	return f, nil
}
