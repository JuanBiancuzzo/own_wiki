package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_MATERIA = "INSERT INTO materias (nombre, codigo, etapa, idCarrera, idPlan, idCuatrimestre, idArchivo) VALUES (?, ?, ?, ?, ?, ?, ?)"
const QUERY_MATERIA_PATH = `SELECT res.id FROM (
	SELECT materias.id, archivos.path FROM archivos INNER JOIN materias ON archivos.id = materias.idArchivo
) AS res WHERE res.path = ?`

const QUERY_PLANES = "SELECT id FROM planesCarrera WHERE nombre = ?"
const INSERTAR_PLAN = "INSERT INTO planesCarrera (nombre) VALUES (?)"

const QUERY_CUATRIMESTRES = "SELECT id FROM cuatrimestreCarrera WHERE anio = ? AND cuatrimestre = ?"
const INSERTAR_CUATRIMESTRE = "INSERT INTO cuatrimestreCarrera (anio, cuatrimestre) VALUES (?, ?)"

const INSERTAR_CORRELATIVAS = "INSERT INTO materiasCorrelativas (idMateria, idCorrelativa) VALUES (?, ?)"

type ParteCuatrimestre string

const (
	CUATRIMESTRE_PRIMERO = "Primero"
	CUATRIMESTRE_SEGUNDO = "Segundo"
)

type Materia struct {
	IdArchivo         *Opcional[int64]
	IdCarrera         *Opcional[int64]
	Nombre            string
	Codigo            string
	Plan              string
	Anio              int
	Cuatri            ParteCuatrimestre
	Etapa             Etapa
	ListaDependencias *l.Lista[Dependencia]
}

func NewMateria(nombre string, codigo string, plan string, repCuatri string, repEtapa string) (*Materia, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else if anio, cuatri, err := ObtenerCuatrimestreParte(repCuatri); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else {
		return &Materia{
			IdArchivo:         NewOpcional[int64](),
			IdCarrera:         NewOpcional[int64](),
			Nombre:            nombre,
			Codigo:            codigo,
			Anio:              anio,
			Cuatri:            cuatri,
			Etapa:             etapa,
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (m *Materia) CrearDependenciaCarrera(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		m.IdCarrera.Asignar(id)
		return m, CumpleAll(m.IdArchivo, m.IdCarrera)
	})
}

func (m *Materia) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		m.IdArchivo.Asignar(id)
		return m, CumpleAll(m.IdArchivo, m.IdCarrera)
	})
}

func (m *Materia) CargarDependencia(dependencia Dependencia) {
	m.ListaDependencias.Push(dependencia)
}

func (m *Materia) Insertar(idPlan int64, idCuatrimestre int64) ([]any, error) {
	if idArchivo, existe := m.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia no tiene todavia el idArchivo")

	} else if idCarrera, existe := m.IdCarrera.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia no tiene todavia el idCarrera")

	} else {
		return []any{m.Nombre, m.Codigo, m.Etapa, idCarrera, idPlan, idCuatrimestre, idArchivo}, nil
	}
}

func (m *Materia) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- fmt.Sprintf("Insertar Materia: %s", m.Nombre)

	if idPlan, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_PLANES, m.Plan) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_PLAN, m.Plan) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry del plan %s con error: %v", m.Plan, err)

	} else if idCuatrimestre, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_CUATRIMESTRES, m.Anio, m.Cuatri) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_CUATRIMESTRE, m.Anio, m.Cuatri) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry del cuatri %s parte de %d con error: %v", m.Cuatri, m.Anio, err)

	} else if datos, err := m.Insertar(idPlan, idCuatrimestre); err != nil {
		return 0, err

	} else {
		return InsertarDirecto(bdd, INSERTAR_MATERIA, datos...)
	}
}

func (m *Materia) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, m.ListaDependencias.Items())
}

func ObtenerCuatrimestreParte(representacionCuatri string) (int, ParteCuatrimestre, error) {
	var anio int
	var cuatriNum int
	var cuatri ParteCuatrimestre

	if _, err := fmt.Sscanf(representacionCuatri, "%dC%d", &anio, &cuatriNum); err != nil {
		return anio, cuatri, fmt.Errorf("el tipo de anio-cuatri (%s) no es uno de los esperados", representacionCuatri)
	}

	switch cuatriNum {
	case 1:
		cuatri = CUATRIMESTRE_PRIMERO
	case 2:
		cuatri = CUATRIMESTRE_SEGUNDO
	default:
		return anio, cuatri, fmt.Errorf("el cuatri dado por %d no es posible representar", cuatriNum)
	}

	return anio, cuatri, nil
}
