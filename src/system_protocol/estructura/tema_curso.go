package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const INSERTAR_TEMA_CURSO = "INSERT INTO temasCurso (nombre, capitulo, parte, tipoCurso, idCurso, idArchivo) VALUES (?, ?, ?, ?, ?, ?)"

type TipoCurso string

const (
	CURSO_ONLINE     = "Online"
	CURSO_PRESENCIAL = "Presencial"
)

type ConstructorTemaCurso struct {
	IdArchivo         *Opcional[int64]
	TipoCurso         TipoCurso
	IdCurso           *Opcional[int64]
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorTemaCurso(nombre string, capitulo string, parte string, tipoCurso TipoCurso) *ConstructorTemaCurso {
	return &ConstructorTemaCurso{
		IdArchivo:         NewOpcional[int64](),
		TipoCurso:         tipoCurso,
		IdCurso:           NewOpcional[int64](),
		Nombre:            nombre,
		Capitulo:          NumeroODefault(capitulo, 1),
		Parte:             NumeroODefault(parte, 0),
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (ctc *ConstructorTemaCurso) CumpleDependencia() (*TemaCurso, bool) {
	if idArchivo, existe := ctc.IdArchivo.Obtener(); !existe {
		return nil, false

	} else if idCurso, existe := ctc.IdCurso.Obtener(); !existe {
		return nil, false

	} else {
		return &TemaCurso{
			IdArchivo:         idArchivo,
			TipoCurso:         ctc.TipoCurso,
			IdCurso:           idCurso,
			Nombre:            ctc.Nombre,
			Capitulo:          ctc.Capitulo,
			Parte:             ctc.Parte,
			ListaDependencias: ctc.ListaDependencias,
		}, true
	}
}

func (ctc *ConstructorTemaCurso) CrearDependenciaCurso(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		ctc.IdCurso.Asignar(id)
		return ctc.CumpleDependencia()
	})
}

func (ctc *ConstructorTemaCurso) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		ctc.IdArchivo.Asignar(id)
		return ctc.CumpleDependencia()
	})
}

func (ctc *ConstructorTemaCurso) CargarDependencia(dependencia Dependencia) {
	ctc.ListaDependencias.Push(dependencia)
}

type TemaCurso struct {
	IdArchivo         int64
	IdCurso           int64
	TipoCurso         TipoCurso
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func (tc *TemaCurso) Insertar() []any {
	return []any{tc.Nombre, tc.Capitulo, tc.Parte, tc.TipoCurso, tc.IdCurso, tc.IdArchivo}
}

func (tc *TemaCurso) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Resumen Curso: %s", tc.Nombre)
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_TEMA_CURSO, tc.Insertar()...) },
	)
}

func (tc *TemaCurso) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, tc.ListaDependencias.Items())
}
