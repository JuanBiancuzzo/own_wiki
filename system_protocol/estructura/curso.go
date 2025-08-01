package estructura

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
	"strconv"
)

const INSERTAR_CURSO = "INSERT INTO cursos (nombre, etapa, anioCurso, idPagina, url, idArchivo) VALUES (?, ?, ?, ?, ?, ?)"

const QUERY_PAGINA_CURSO = "SELECT id FROM paginasCursos WHERE nombrePagina = ?"
const INSERTAR_PAGINA_CURSO = "INSERT INTO paginasCursos (nombrePagina) VALUES (?)"

const INSERTAR_PROFESOR_CURSO = "INSERT INTO profesoresCurso (idCurso, tipoCurso, idPersona) VALUES (?, ?, ?)"

type ConstructorCurso struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	NombrePagina      string
	Url               string
	Profesores        []*Persona
	ListaDependencias *l.Lista[Dependencia]
}

func NewConstructorCurso(nombre string, repEtapa string, repAnioCurso string, nombrePagina string, url string, profesores []*Persona) (*ConstructorCurso, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear curso con error: %v", err)

	} else if anioCurso, err := strconv.Atoi(repAnioCurso); err != nil {
		return nil, fmt.Errorf("error al crear curso al obtener el anio, con error: %v", err)

	} else {
		return &ConstructorCurso{
			Nombre:            nombre,
			Etapa:             etapa,
			AnioCurso:         anioCurso,
			NombrePagina:      nombrePagina,
			Url:               url,
			Profesores:        profesores,
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (cc *ConstructorCurso) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		return &Curso{
			Nombre:            cc.Nombre,
			Etapa:             cc.Etapa,
			AnioCurso:         cc.AnioCurso,
			NombrePagina:      cc.NombrePagina,
			Url:               cc.Url,
			Profesores:        cc.Profesores,
			IdArchivo:         id,
			ListaDependencias: cc.ListaDependencias,
		}, true
	})
}

func (cc *ConstructorCurso) CargarDependencia(dependencia Dependencia) {
	cc.ListaDependencias.Push(dependencia)
}

type Curso struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	NombrePagina      string
	Url               string
	IdArchivo         int64
	Profesores        []*Persona
	ListaDependencias *l.Lista[Dependencia]
}

func (c *Curso) Insertar(idPagina int64) []any {
	return []any{c.Nombre, c.Etapa, c.AnioCurso, idPagina, c.Url, c.IdArchivo}
}

func (c *Curso) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Curso"

	var idCurso int64

	if idPagina, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_PAGINA_CURSO, c.NombrePagina) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_PAGINA_CURSO, c.NombrePagina) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry del nombre de la pagina: '%s' con error: %v", c.NombrePagina, err)

	} else if idCurso, err = Insertar(func() (sql.Result, error) {
		return bdd.Exec(INSERTAR_CURSO, c.Insertar(idPagina)...)
	}); err != nil {
		return 0, fmt.Errorf("error al insertar un curso, con error: %v", err)
	}

	for _, profesor := range c.Profesores {
		if idAutor, err := ObtenerOInsertar(
			func() *sql.Row { return bdd.QueryRow(QUERY_PERSONAS, profesor.Insertar()...) },
			func() (sql.Result, error) { return bdd.Exec(INSERTAR_PERSONA, profesor.Insertar()...) },
		); err != nil {
			canal <- fmt.Sprintf("error al hacer una querry del profesor: %s %s con error: %v", profesor.Nombre, profesor.Apellido, err)

		} else if _, err := bdd.Exec(INSERTAR_PROFESOR_CURSO, idCurso, CURSO_ONLINE, idAutor); err != nil {
			canal <- fmt.Sprintf("error al insertar par idCurso-idProfesor, con error: %v", err)
		}
	}

	return idCurso, nil
}

func (c *Curso) ResolverDependencias(id int64) []Cargable {
	return ResolverDependencias(id, c.ListaDependencias.Items())
}
