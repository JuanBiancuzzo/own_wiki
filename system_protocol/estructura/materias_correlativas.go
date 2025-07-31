package estructura

import (
	"database/sql"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_MATERIA_CORRELATIVA = "INSERT INTO materiasCorrelativas (tipoMateria, idMateria, tipoCorrelativa, idCorrelativa) VALUES (?, ?, ?, ?)"

type TipoMateria string

const (
	MATERIA_REAL        = "Materia"
	MATERIA_EQUIVALENTE = "Equivalente"
)

type ConstructorMateriasCorrelativas struct {
	IdMateria         Opcional[int64]
	TipoMateria       TipoMateria
	IdCorrelativa     Opcional[int64]
	PathCorrelativa   string
	TipoCorrelativa   TipoMateria
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorMateriasCorrelativas(tipoMateria TipoMateria, pathCorrelativa string) *ConstructorMateriasCorrelativas {
	return &ConstructorMateriasCorrelativas{
		IdMateria:         NewOpcional[int64](),
		TipoMateria:       tipoMateria,
		IdCorrelativa:     NewOpcional[int64](),
		PathCorrelativa:   pathCorrelativa,
		TipoCorrelativa:   "",
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (cmc *ConstructorMateriasCorrelativas) CumpleDependencia() (*MateriasCorrelativas, bool) {
	if cmc.IdMateria.Esta && cmc.IdCorrelativa.Esta {
		return &MateriasCorrelativas{
			IdMateria:         cmc.IdMateria.Valor,
			TipoArchivo:       cmc.TipoMateria,
			IdCorrelativa:     cmc.IdCorrelativa.Valor,
			TipoCorrelativa:   cmc.TipoCorrelativa,
			ListaDependencias: cmc.ListaDependencias,
		}, true
	}

	return nil, false
}

func (cmc *ConstructorMateriasCorrelativas) CumpleDependenciaCorrelativa(id int64) (Cargable, bool) {
	cmc.IdCorrelativa.Asignar(id)
	return cmc.CumpleDependencia()
}

func (cmc *ConstructorMateriasCorrelativas) CumpleDependenciaMateria(id int64) (Cargable, bool) {
	cmc.IdMateria.Asignar(id)
	return cmc.CumpleDependencia()
}

func (cmc *ConstructorMateriasCorrelativas) CargarDependencia(dependencia Dependencia) {
	cmc.ListaDependencias.Push(dependencia)
}

type MateriasCorrelativas struct {
	IdMateria         int64
	TipoArchivo       TipoMateria
	IdCorrelativa     int64
	TipoCorrelativa   TipoMateria
	ListaDependencias *l.Lista[Dependencia]
}

func (mc *MateriasCorrelativas) Insertar() []any {
	return []any{mc.TipoArchivo, mc.IdMateria, mc.TipoCorrelativa, mc.IdCorrelativa}
}

func (mc *MateriasCorrelativas) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Materia Correlativas"
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_MATERIA_CORRELATIVA, mc.Insertar()...) },
	)
}

func (mc *MateriasCorrelativas) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, mc.ListaDependencias)
}
