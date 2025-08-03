package datos

import (
	"database/sql"
	"fmt"
	u "own_wiki/system_protocol/utilidades"
)

const INSERTAR_MATERIA_EQUIVALENTES = "INSERT INTO materiasEquivalentes (nombre, codigo, idCarrera, idMateria, idArchivo) VALUES (?, ?, ?, ?, ?)"

type MateriaEquivalente struct {
	IdArchivo         *Opcional[int64]
	IdCarrera         *Opcional[int64]
	IdMateria         *Opcional[int64]
	Nombre            string
	Codigo            string
	ListaDependencias *u.Lista[Dependencia]
}

func NewMateriaEquivalente(nombre string, codigo string) *MateriaEquivalente {
	return &MateriaEquivalente{
		IdArchivo:         NewOpcional[int64](),
		IdCarrera:         NewOpcional[int64](),
		IdMateria:         NewOpcional[int64](),
		Nombre:            nombre,
		Codigo:            codigo,
		ListaDependencias: u.NewLista[Dependencia](),
	}
}

func (me *MateriaEquivalente) CrearDependenciaCarrera(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		me.IdCarrera.Asignar(id)
		return me, CumpleAll(me.IdArchivo, me.IdCarrera, me.IdMateria)
	})
}
func (me *MateriaEquivalente) CrearDependenciaMateria(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		me.IdMateria.Asignar(id)
		return me, CumpleAll(me.IdArchivo, me.IdCarrera, me.IdMateria)
	})
}

func (me *MateriaEquivalente) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		me.IdArchivo.Asignar(id)
		return me, CumpleAll(me.IdArchivo, me.IdCarrera, me.IdMateria)
	})
}

func (me *MateriaEquivalente) CargarDependencia(dependencia Dependencia) {
	me.ListaDependencias.Push(dependencia)
}

func (me *MateriaEquivalente) Insertar() ([]any, error) {
	if idArchivo, existe := me.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia equivalente no tiene todavia el idArchivo")

	} else if idMateria, existe := me.IdMateria.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia equivalente no tiene todavia el idMateria")

	} else if idCarrera, existe := me.IdCarrera.Obtener(); !existe {
		return []any{}, fmt.Errorf("materia equivalente no tiene todavia el idCarrera")

	} else {
		return []any{me.Nombre, me.Codigo, idCarrera, idMateria, idArchivo}, nil
	}
}

func (me *MateriaEquivalente) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- fmt.Sprintf("Insertar Materia Equivalentes: %s", me.Nombre)
	if datos, err := me.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_MATERIA_EQUIVALENTES, datos...)
	}
}

func (me *MateriaEquivalente) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, me.ListaDependencias.Items())
}
