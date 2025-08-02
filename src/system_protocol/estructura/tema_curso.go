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

type TemaCurso struct {
	IdArchivo         *Opcional[int64]
	TipoCurso         TipoCurso
	IdCurso           *Opcional[int64]
	Nombre            string
	Capitulo          int
	Parte             int
	ListaDependencias *l.Lista[Dependencia]
}

func NewTemaCurso(nombre string, capitulo string, parte string, tipoCurso TipoCurso) *TemaCurso {
	return &TemaCurso{
		IdArchivo:         NewOpcional[int64](),
		TipoCurso:         tipoCurso,
		IdCurso:           NewOpcional[int64](),
		Nombre:            nombre,
		Capitulo:          NumeroODefault(capitulo, 1),
		Parte:             NumeroODefault(parte, 0),
		ListaDependencias: l.NewLista[Dependencia](),
	}
}

func (tc *TemaCurso) CrearDependenciaCurso(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		tc.IdCurso.Asignar(id)
		return tc, CumpleAll(tc.IdArchivo, tc.IdCurso)
	})
}

func (tc *TemaCurso) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		tc.IdArchivo.Asignar(id)
		return tc, CumpleAll(tc.IdArchivo, tc.IdCurso)
	})
}

func (tc *TemaCurso) CargarDependencia(dependencia Dependencia) {
	tc.ListaDependencias.Push(dependencia)
}

func (tc *TemaCurso) Insertar() ([]any, error) {
	if idArchivo, existe := tc.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("tema curso no tiene todavia el idArchivo")

	} else if idCurso, existe := tc.IdCurso.Obtener(); !existe {
		return []any{}, fmt.Errorf("tema curso no tiene todavia el idCurso")

	} else {
		return []any{tc.Nombre, tc.Capitulo, tc.Parte, tc.TipoCurso, idCurso, idArchivo}, nil
	}
}

func (tc *TemaCurso) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- fmt.Sprintf("Insertar Resumen Curso: %s", tc.Nombre)
	if datos, err := tc.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_TEMA_CURSO, datos...)
	}
}

func (tc *TemaCurso) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, tc.ListaDependencias.Items())
}
