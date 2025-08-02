package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_MATERIA_EQUIVALENTES = "INSERT INTO materiasEquivalentes (nombre, codigo, idCarrera, idMateria, idArchivo) VALUES (?, ?, ?, ?, ?)"

type ConstructorMateriaEquivalente struct {
	IdArchivo         *Opcional[int64]
	idCarrera         *Opcional[int64]
	IdMateria         *Opcional[int64]
	Nombre            string
	Codigo            string
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorMateriaEquivalente(nombre string, codigo string) *ConstructorMateriaEquivalente {
	return &ConstructorMateriaEquivalente{
		IdArchivo:         NewOpcional[int64](),
		idCarrera:         NewOpcional[int64](),
		IdMateria:         NewOpcional[int64](),
		Nombre:            nombre,
		Codigo:            codigo,
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (cme *ConstructorMateriaEquivalente) CumpleDependencia() (*MateriaEquivalente, bool) {
	if idArchivo, existe := cme.IdArchivo.Obtener(); !existe {
		return nil, false

	} else if idMateria, existe := cme.IdMateria.Obtener(); !existe {
		return nil, false

	} else if idCarrera, existe := cme.idCarrera.Obtener(); !existe {
		return nil, false

	} else {
		return &MateriaEquivalente{
			IdArchivo:         idArchivo,
			IdCarrera:         idCarrera,
			IdMateria:         idMateria,
			Nombre:            cme.Nombre,
			Codigo:            cme.Codigo,
			ListaDependencias: cme.ListaDependencias,
		}, true
	}
}

func (cme *ConstructorMateriaEquivalente) CrearDependenciaCarrera(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		cme.idCarrera.Asignar(id)
		return cme.CumpleDependencia()
	})
}
func (cme *ConstructorMateriaEquivalente) CrearDependenciaMateria(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		cme.IdMateria.Asignar(id)
		return cme.CumpleDependencia()
	})
}

func (cme *ConstructorMateriaEquivalente) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		cme.IdArchivo.Asignar(id)
		return cme.CumpleDependencia()
	})
}

func (cme *ConstructorMateriaEquivalente) CargarDependencia(dependencia Dependencia) {
	cme.ListaDependencias.Push(dependencia)
}

type MateriaEquivalente struct {
	IdArchivo         int64
	IdMateria         int64
	IdCarrera         int64
	Nombre            string
	Codigo            string
	ListaDependencias *l.Lista[Dependencia]
}

func (me *MateriaEquivalente) Insertar() []any {
	return []any{me.Nombre, me.Codigo, me.IdCarrera, me.IdMateria, me.IdArchivo}
}

func (me *MateriaEquivalente) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Materia Equivalentes: %s", me.Nombre)
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_MATERIA_EQUIVALENTES, me.Insertar()...) },
	)
}

func (me *MateriaEquivalente) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, me.ListaDependencias.Items())
}
