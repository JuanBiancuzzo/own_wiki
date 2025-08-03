package fs

import (
	"database/sql"
	"fmt"
	l "own_wiki/system_protocol/listas"
)

const (
	QUERY_CURSOS_LS      = "SELECT nombre FROM cursos UNION ALL SELECT nombre FROM cursosPresencial"
	QUERY_TEMAS_CURSO_LS = "SELECT nombre FROM temasCurso WHERE idCurso = %d"
	QUERY_NOTA_CURSO_LS  = `SELECT DISTINCT notas.nombre FROM notas INNER JOIN (
		SELECT idNota, idVinculo FROM notasVinculo WHERE tipoVinculo = "Curso"
	) AS vinculo ON notas.id = vinculo.idNota WHERE vinculo.idVinculo = %d`
)

const (
	QUERY_OBTENER_CURSO = `SELECT id FROM cursos WHERE nombre = '%s'
		UNION ALL
	SELECT idMateria AS id FROM materiasEquivalentes WHERE nombre = '%s'
	`
	QUERY_OBTENER_TEMA_CURSO = "SELECT id FROM temasCurso WHERE idCurso = %d AND nombre = '%s'"
)

type TipoCurso byte

const (
	TCC_GENERAL = iota
	TCC_DENTRO_CURSO
	TCC_DENTRO_TEMA
)

type Cursos struct {
	Bdd  *sql.DB
	Tipo TipoCurso
	Path *l.Pila[int64]
}

func NewCursos(bdd *sql.DB) *Cursos {
	return &Cursos{
		Bdd:  bdd,
		Tipo: TCC_GENERAL,
		Path: l.NewPila[int64](),
	}
}

func (c *Cursos) Ls() ([]string, error) {
	var query string
	switch c.Tipo {
	case TCC_GENERAL:
		query = QUERY_CURSOS_LS
	case TCC_DENTRO_CURSO:
		if idCurso, err := c.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_TEMAS_CURSO_LS, idCurso)
		}
	case TCC_DENTRO_TEMA:
		if idTema, err := c.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_NOTA_CURSO_LS, idTema)
		}
	}

	if rows, err := c.Bdd.Query(query); err != nil {
		return []string{}, fmt.Errorf("se obtuvo un error en cursos, al hacer query, dando el error: %v", err)

	} else {
		columnas := l.NewLista[string]()
		defer rows.Close()
		for rows.Next() {
			var nombre string
			_ = rows.Scan(&nombre)
			columnas.Push(nombre)
		}

		return columnas.Items(), nil
	}
}

func (c *Cursos) Cd(subpath string, cache *Cache) (Subpath, error) {
	if subpath == ".." {
		return c.RutinaAtras(cache)
	}

	query := ""
	switch c.Tipo {
	case TCC_GENERAL:
		query = fmt.Sprintf(QUERY_OBTENER_CURSO, subpath, subpath)
	case TCC_DENTRO_CURSO:
		if idCurso, err := c.Path.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTENER_TEMA_CURSO, idCurso, subpath)
		}
	case TCC_DENTRO_TEMA:
		return c, fmt.Errorf("ya se esta viendo todos los archivos, no hay subcarpetas")
	}

	if query == "" {
		return c, fmt.Errorf("hubo un error en la query, y esta vacia")
	}

	fila := c.Bdd.QueryRow(query)
	var id int64
	if err := fila.Scan(&id); err != nil {
		return c, fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	switch c.Tipo {
	case TCC_GENERAL:
		c.Tipo = TCC_DENTRO_CURSO
	case TCC_DENTRO_CURSO:
		c.Tipo = TCC_DENTRO_TEMA
	}

	c.Path.Apilar(id)
	return c, nil
}

func (c *Cursos) RutinaAtras(cache *Cache) (Subpath, error) {
	_, _ = c.Path.Desapilar()

	switch c.Tipo {
	case TCC_GENERAL:
		return cache.ObtenerSubpath(PD_ROOT)

	case TCC_DENTRO_CURSO:
		c.Tipo = TCC_GENERAL
	case TCC_DENTRO_TEMA:
		c.Tipo = TCC_DENTRO_CURSO
	}

	return c, nil
}
