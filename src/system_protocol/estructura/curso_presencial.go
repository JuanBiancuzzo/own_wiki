package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
	"strconv"
)

const INSERTAR_CURSO_PRESENCIAL = "INSERT INTO cursosPresencial (nombre, etapa, anioCurso, idArchivo) VALUES (?, ?, ?, ?)"

type ConstructorCursoPresencial struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	Profesores        []*Persona
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorCursoPresencial(nombre string, repEtapa string, repAnioCurso string, profesores []*Persona) (*ConstructorCursoPresencial, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear curso presencial con error: %v", err)

	} else if anioCurso, err := strconv.Atoi(repAnioCurso); err != nil {
		return nil, fmt.Errorf("error al crear curso presencial al obtener el anio, con error: %v", err)

	} else {
		return &ConstructorCursoPresencial{
			Nombre:            nombre,
			Etapa:             etapa,
			AnioCurso:         anioCurso,
			Profesores:        profesores,
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (ccp *ConstructorCursoPresencial) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		return &CursoPresencial{
			Nombre:            ccp.Nombre,
			Etapa:             ccp.Etapa,
			AnioCurso:         ccp.AnioCurso,
			Profesores:        ccp.Profesores,
			IdArchivo:         id,
			ListaDependencias: ccp.ListaDependencias,
		}, true
	})
}

func (ccp *ConstructorCursoPresencial) CargarDependencia(dependencia Dependencia) {
	ccp.ListaDependencias.Push(dependencia)
}

type CursoPresencial struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	IdArchivo         int64
	Profesores        []*Persona
	ListaDependencias *l.Lista[Dependencia]
}

func (c *CursoPresencial) Insertar() []any {
	return []any{c.Nombre, c.Etapa, c.AnioCurso, c.IdArchivo}
}

func (c *CursoPresencial) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Curso Presencial"

	idCursoPresencial, err := Insertar(func() (sql.Result, error) {
		return bdd.Exec(INSERTAR_CURSO_PRESENCIAL, c.Insertar()...)
	})
	if err != nil {
		return 0, fmt.Errorf("error al insertar un curso, con error: %v", err)
	}

	for _, profesor := range c.Profesores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, profesor.Insertar()...) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, profesor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del profesor: %s %s con error: %v", profesor.Nombre, profesor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_PROFESOR_CURSO, idCursoPresencial, CURSO_PRESENCIAL, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idCurso-idProfesor, con error: %v", err)
		}
	}

	return idCursoPresencial, nil
}

func (c *CursoPresencial) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, c.ListaDependencias.Items())
}
