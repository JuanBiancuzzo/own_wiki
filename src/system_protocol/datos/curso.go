package datos

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

type Curso struct {
	Nombre            string
	Etapa             Etapa
	AnioCurso         int
	NombrePagina      string
	Url               string
	IdArchivo         *Opcional[int64]
	Profesores        []*Persona
	ListaDependencias *l.Lista[Dependencia]
}

func NewCurso(nombre string, repEtapa string, repAnioCurso string, nombrePagina string, url string, profesores []*Persona) (*Curso, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear curso con error: %v", err)

	} else if anioCurso, err := strconv.Atoi(repAnioCurso); err != nil {
		return nil, fmt.Errorf("error al crear curso al obtener el anio, con error: %v", err)

	} else {
		return &Curso{
			Nombre:            nombre,
			Etapa:             etapa,
			AnioCurso:         anioCurso,
			NombrePagina:      nombrePagina,
			Url:               url,
			IdArchivo:         NewOpcional[int64](),
			Profesores:        profesores,
			ListaDependencias: l.NewLista[Dependencia](),
		}, nil
	}
}

func (c *Curso) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		c.IdArchivo.Asignar(id)
		return c, true
	})
}

func (c *Curso) CargarDependencia(dependencia Dependencia) {
	c.ListaDependencias.Push(dependencia)
}

func (c *Curso) Insertar(idPagina int64) ([]any, error) {
	if idArchivo, existe := c.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("curso no tiene todavia el idArchivo")
	} else {
		return []any{c.Nombre, c.Etapa, c.AnioCurso, idPagina, c.Url, idArchivo}, nil
	}
}

func (c *Curso) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- "Insertar Curso"

	var idCurso int64

	if idPagina, err := ObtenerOInsertar(
		func() *sql.Row { return bdd.QueryRow(QUERY_PAGINA_CURSO, c.NombrePagina) },
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_PAGINA_CURSO, c.NombrePagina) },
	); err != nil {
		return 0, fmt.Errorf("error al hacer una querry del nombre de la pagina: '%s' con error: %v", c.NombrePagina, err)

	} else if datos, err := c.Insertar(idPagina); err != nil {
		return 0, err

	} else if idCurso, err = InsertarDirecto(bdd, INSERTAR_CURSO, datos...); err != nil {
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
