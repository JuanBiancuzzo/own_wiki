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

type Opcional[T any] struct {
	Valor T
	Esta  bool
}

func NewOpcional[T any]() Opcional[T] {
	var valor T
	return Opcional[T]{
		Valor: valor,
		Esta:  false,
	}
}

func (o Opcional[T]) Asignar(valor T) {
	o.Valor = valor
	o.Esta = true
}

type ConstructorMateria struct {
	IdArchivo         Opcional[int64]
	IdCarrera         Opcional[int64]
	Nombre            string
	Codigo            string
	Plan              string
	Anio              int
	Cuatri            ParteCuatrimestre
	Etapa             Etapa
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorMateria(nombre string, codigo string, plan string, repCuatri string, repEtapa string) (*ConstructorMateria, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else if anio, cuatri, err := ObtenerCuatrimestreParte(repCuatri); err != nil {
		return nil, fmt.Errorf("error al crear materia con error: %v", err)

	} else {
		return &ConstructorMateria{
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

func (cm *ConstructorMateria) CumpleDependencia() (*Materia, bool) {
	if cm.IdArchivo.Esta && cm.IdCarrera.Esta {
		return &Materia{
			Nombre:            cm.Nombre,
			Codigo:            cm.Codigo,
			Plan:              cm.Plan,
			Anio:              cm.Anio,
			Cuatri:            cm.Cuatri,
			Etapa:             cm.Etapa,
			IdCarrera:         cm.IdCarrera.Valor,
			IdArchivo:         cm.IdArchivo.Valor,
			ListaDependencias: cm.ListaDependencias,
		}, true
	}
	return nil, false
}

func (cm *ConstructorMateria) CumpleDependenciaCarrera(id int64) (Cargable, bool) {
	cm.IdCarrera.Asignar(id)
	return cm.CumpleDependencia()
}

func (cm *ConstructorMateria) CumpleDependenciaArchivo(id int64) (Cargable, bool) {
	cm.IdArchivo.Asignar(id)
	return cm.CumpleDependencia()
}

func (cm *ConstructorMateria) CargarDependencia(dependencia Dependencia) {
	cm.ListaDependencias.Push(dependencia)
}

type Materia struct {
	Nombre            string
	Codigo            string
	Plan              string
	Anio              int
	Cuatri            ParteCuatrimestre
	Etapa             Etapa
	IdCarrera         int64
	IdArchivo         int64
	ListaDependencias *l.Lista[Dependencia]
}

func (m *Materia) Insertar(idPlan int64, idCuatrimestre int64) []any {
	return []any{m.Nombre, m.Codigo, m.Etapa, m.IdCarrera, idPlan, idCuatrimestre, m.IdArchivo}
}

func (m *Materia) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Materia: %s", m.Nombre)

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

	} else {
		return Insertar(
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_MATERIA, m.Insertar(idPlan, idCuatrimestre)...) },
		)
	}
}

func (m *Materia) ResolverDependencias(id int64) []Cargable {
	cantidadCumple := 0
	cargables := make([]Cargable, m.ListaDependencias.Largo)

	for cumpleDependencia := range m.ListaDependencias.Iterar {
		if cargable, cumple := cumpleDependencia(id); cumple {
			cargables[cantidadCumple] = cargable
			cantidadCumple++
		}
	}

	return cargables[:cantidadCumple]
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
