package datos

import (
	"database/sql"
	"fmt"
	u "own_wiki/system_protocol/utilidades"
)

const INSERTAR_MATERIA_CORRELATIVA = "INSERT INTO materiasCorrelativas (tipoMateria, idMateria, tipoCorrelativa, idCorrelativa) VALUES (?, ?, ?, ?)"

type TipoMateria string

const (
	MATERIA_REAL        = "Materia"
	MATERIA_EQUIVALENTE = "Equivalente"
)

type MateriasCorrelativas struct {
	IdMateria         *Opcional[int64]
	TipoMateria       TipoMateria
	IdCorrelativa     *Opcional[int64]
	TipoCorrelativa   TipoMateria
	ListaDependencias *u.Lista[Dependencia]
}

func NewMateriasCorrelativas(tipoMateria TipoMateria, tipoCorrelativa TipoMateria) *MateriasCorrelativas {
	return &MateriasCorrelativas{
		IdMateria:         NewOpcional[int64](),
		TipoMateria:       tipoMateria,
		IdCorrelativa:     NewOpcional[int64](),
		TipoCorrelativa:   tipoCorrelativa,
		ListaDependencias: u.NewLista[Dependencia](),
	}
}

func (mc *MateriasCorrelativas) CrearDependenciaCorrelativa(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		mc.IdCorrelativa.Asignar(id)
		return mc, CumpleAll(mc.IdMateria, mc.IdCorrelativa)
	})
}

func (mc *MateriasCorrelativas) CrearDependenciaMateria(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		mc.IdMateria.Asignar(id)
		return mc, CumpleAll(mc.IdMateria, mc.IdCorrelativa)
	})
}

func (mc *MateriasCorrelativas) CargarDependencia(dependencia Dependencia) {
	mc.ListaDependencias.Push(dependencia)
}

func (mc *MateriasCorrelativas) Insertar() ([]any, error) {
	if idMateria, existe := mc.IdMateria.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia correlativa no tiene todavia el idArchivo")

	} else if idCorrelativa, existe := mc.IdCorrelativa.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia correlativa no tiene todavia el idCorrelativa")

	} else {
		return []any{mc.TipoMateria, idMateria, mc.TipoCorrelativa, idCorrelativa}, nil
	}
}

func (mc *MateriasCorrelativas) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- "Insertar Materia Correlativas"
	if datos, err := mc.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_MATERIA_CORRELATIVA, datos...)
	}
}

func (mc *MateriasCorrelativas) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, mc.ListaDependencias.Items())
}
