package fs

import (
	"database/sql"
	"fmt"

	l "own_wiki/system_protocol/utilidades"
)

const (
	QUERY_CARRERAS_LS = "SELECT nombre FROM carreras"
	QUERY_MATERIAS_LS = `SELECT materiasGlobal.nombre FROM (
		SELECT idCarrera, nombre FROM materias UNION ALL SELECT idCarrera, nombre FROM materiasEquivalentes
	) AS materiasGlobal WHERE materiasGlobal.idCarrera = %d`
	QUERY_TEMAS_MATERIA_LS = "SELECT nombre FROM temasMateria WHERE idMateria = %d"
	QUERY_NOTA_MATERIA_LS  = `SELECT DISTINCT notas.nombre FROM notas INNER JOIN (
		SELECT idNota, idVinculo FROM notasVinculo WHERE tipoVinculo = "Facultad"
	) AS vinculo ON notas.id = vinculo.idNota WHERE vinculo.idVinculo = %d`
)

const (
	QUERY_OBTENER_CARRERA = "SELECT id, nombre FROM carreras WHERE nombre = '%s'"
	QUERY_OBTENER_MATERIA = `SELECT id, nombre FROM materias WHERE idCarrera = %d AND nombre = '%s'
		UNION ALL
	SELECT idMateria AS id, nombre FROM materiasEquivalentes WHERE idCarrera = %d AND nombre = '%s'
	`
	QUERY_OBTNER_TEMA_MATERIA = "SELECT id, nombre FROM temasMateria WHERE idMateria = %d AND nombre = '%s'"
)

type TipoFacultad byte

const (
	TF_FACULTAD = iota
	TF_DENTRO_CARRERA
	TF_DENTRO_MATERIA
	TF_DENTRO_TEMA
)

type Facultad struct {
	Bdd  *sql.DB
	Tipo TipoFacultad
	Path *l.Pila[int64]
}

func NewFacultad(bdd *sql.DB) *Facultad {
	return &Facultad{
		Bdd:  bdd,
		Tipo: TF_FACULTAD,
		Path: l.NewPila[int64](),
	}
}

func (f *Facultad) Ls() ([]string, error) {
	var query string
	switch f.Tipo {
	case TF_FACULTAD:
		query = QUERY_CARRERAS_LS
	case TF_DENTRO_CARRERA:
		if idCarrera, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_MATERIAS_LS, idCarrera)
		}
	case TF_DENTRO_MATERIA:
		if idMateria, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_TEMAS_MATERIA_LS, idMateria)
		}
	case TF_DENTRO_TEMA:
		if idTema, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_NOTA_MATERIA_LS, idTema)
		}

	}

	if rows, err := f.Bdd.Query(query); err != nil {
		return []string{}, fmt.Errorf("se obtuvo un error en facultad, al hacer query, dando el error: %v", err)

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

func (f *Facultad) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		return f.RutinaAtras(cache)
	}

	query := ""
	switch f.Tipo {
	case TF_FACULTAD:
		query = fmt.Sprintf(QUERY_OBTENER_CARRERA, subpath)
	case TF_DENTRO_CARRERA:
		if idCarrera, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTENER_MATERIA, idCarrera, subpath, idCarrera, subpath)
		}
	case TF_DENTRO_MATERIA:
		if idMateria, err := f.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTNER_TEMA_MATERIA, idMateria, subpath)
		}
	case TF_DENTRO_TEMA:
		return f, fmt.Errorf("ya se esta viendo todos los archivos, no hay subcarpetas")
	}

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
	case TF_FACULTAD:
		f.Tipo = TF_DENTRO_CARRERA
	case TF_DENTRO_CARRERA:
		f.Tipo = TF_DENTRO_MATERIA
	case TF_DENTRO_MATERIA:
		f.Tipo = TF_DENTRO_TEMA
	}

	f.Path.Apilar(id)
	return f, nil
}

func (f *Facultad) RutinaAtras(cache *Cache) (Subpath, error) {
	_, _ = f.Path.Desapilar()

	switch f.Tipo {
	case TF_FACULTAD:
		return cache.ObtenerSubpath(PD_ROOT)

	case TF_DENTRO_CARRERA:
		f.Tipo = TF_FACULTAD
	case TF_DENTRO_MATERIA:
		f.Tipo = TF_DENTRO_CARRERA
	case TF_DENTRO_TEMA:
		f.Tipo = TF_DENTRO_MATERIA
	}

	return f, nil
}
