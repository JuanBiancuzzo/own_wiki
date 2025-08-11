package datos

/*

import (
	"database/sql"
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	u "own_wiki/system_protocol/utilidades"
	"strconv"
)

const INSERTAR_CURSO_PRESENCIAL = "INSERT INTO cursosPresencial (nombre, etapa, anioCurso, idArchivo) VALUES (?, ?, ?, ?)"

type CursoPresencial struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	IdArchivo         *u.Opcional[int64]
	Profesores        []*Persona
	ListaDependencias *u.Lista[Dependencia]
}

func NewCursoPresencial(nombre string, repEtapa string, repAnioCurso string, profesores []*Persona) (*CursoPresencial, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear curso presencial con error: %v", err)

	} else if anioCurso, err := strconv.Atoi(repAnioCurso); err != nil {
		return nil, fmt.Errorf("error al crear curso presencial al obtener el anio, con error: %v", err)

	} else {
		return &CursoPresencial{
			Nombre:            nombre,
			Etapa:             etapa,
			AnioCurso:         anioCurso,
			IdArchivo:         u.NewOpcional[int64](),
			Profesores:        profesores,
			ListaDependencias: u.NewLista[Dependencia](),
		}, nil
	}
}

func (cp *CursoPresencial) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		cp.IdArchivo.Asignar(id)
		return cp, true
	})
}

func (cp *CursoPresencial) CargarDependencia(dependencia Dependencia) {
	cp.ListaDependencias.Push(dependencia)
}

func (cp *CursoPresencial) Insertar() ([]any, error) {
	if idArchivo, existe := cp.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("curso presencial no tiene todavia el idArchivo")
	} else {
		return []any{cp.Nombre, cp.Etapa, cp.AnioCurso, idArchivo}, nil
	}
}

func (cp *CursoPresencial) CargarDatos(bdd *b.Bdd, canal chan string) (int64, error) {
	// canal <- "Insertar Curso Presencial"

	var idCursoPresencial int64
	if datos, err := cp.Insertar(); err != nil {
		return 0, err
	} else if idCursoPresencial, err = InsertarDirecto(bdd, INSERTAR_CURSO_PRESENCIAL, datos...); err != nil {
		return 0, fmt.Errorf("error al insertar un curso, con error: %v", err)
	}

	for _, profesor := range cp.Profesores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.MySQL.QueryRow(QUERY_PERSONAS, profesor.Insertar()...) },
			func() (sql.Result, error) { return bdd.MySQL.Exec(INSERTAR_PERSONA, profesor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del profesor: %s %s con error: %v", profesor.Nombre, profesor.Apellido, err)

		} else if _, err := InsertarDirecto(bdd, INSERTAR_PROFESOR_CURSO, idCursoPresencial, CURSO_PRESENCIAL, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idCurso-idProfesor, con error: %v", err)
		}
	}

	return idCursoPresencial, nil
}

func (cp *CursoPresencial) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, cp.ListaDependencias.Items())
}
*/
