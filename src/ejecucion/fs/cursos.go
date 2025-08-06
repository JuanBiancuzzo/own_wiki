package fs

import (
	"database/sql"
	"fmt"
	u "own_wiki/system_protocol/utilidades"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	QUERY_CURSOS_LS      = "SELECT nombre FROM cursos UNION ALL SELECT nombre FROM cursosPresencial"
	QUERY_TEMAS_CURSO_LS = "SELECT nombre FROM temasCurso WHERE idCurso = %d ORDER BY capitulo, parte"
	QUERY_NOTA_CURSO_LS  = `SELECT DISTINCT notas.nombre FROM notas INNER JOIN (
		SELECT idNota, idVinculo FROM notasVinculo WHERE tipoVinculo = "Curso"
	) AS vinculo ON notas.id = vinculo.idNota WHERE vinculo.idVinculo = %d`
)

const (
	QUERY_OBTENER_CURSO = `SELECT id, nombre FROM cursos WHERE nombre = '%s'
		UNION ALL
	SELECT idMateria, nombre AS id FROM materiasEquivalentes WHERE nombre = '%s'
	`
	QUERY_OBTENER_TEMA_CURSO = "SELECT id, nombre FROM temasCurso WHERE idCurso = %d AND nombre = '%s'"
)

type TipoCurso byte

const (
	TCC_GENERAL = iota
	TCC_DENTRO_CURSO
	TCC_DENTRO_TEMA
)

type Cursos struct {
	Bdd    *sql.DB
	Tipo   TipoCurso
	Indice *u.Pila[int64]
	Path   *u.Pila[string]
}

func NewCursos(bdd *sql.DB) *Cursos {
	return &Cursos{
		Bdd:    bdd,
		Tipo:   TCC_GENERAL,
		Indice: u.NewPila[int64](),
		Path:   u.NewPila[string](),
	}
}

func GenerarRutaCursos(e *echo.Echo, bdd *sql.DB) {
	cursos := NewCursos(bdd)

	e.GET("/Cursos", func(c echo.Context) error {
		data := cursos.DecidirSiguientePath(c)
		return c.Render(200, "cursos", data)
	})
}

func (c *Cursos) DecidirSiguientePath(ec echo.Context) Data {
	path := strings.TrimSpace(ec.QueryParam("path"))
	errCd := c.Cd(path)
	data, errLs := c.Ls()

	if errCd != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Cd tuvo el error: %v", errCd), "/Cursos"))
	}

	if errLs != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Ls tuvo el error: %v", errLs), "/Cursos"))
	}

	return data
}

func (c *Cursos) Ls() (Data, error) {
	var data Data
	opciones := u.NewLista[Opcion]()
	returnPath := "/Cursos?path=.."

	var query string
	switch c.Tipo {
	case TCC_GENERAL:
		query = QUERY_CURSOS_LS
		returnPath = "/Root"
	case TCC_DENTRO_CURSO:
		if idCurso, err := c.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_TEMAS_CURSO_LS, idCurso)
		}
	case TCC_DENTRO_TEMA:
		if idTema, err := c.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_NOTA_CURSO_LS, idTema)
		}
	}

	if rows, err := c.Bdd.Query(query); err != nil {
		return data, fmt.Errorf("se obtuvo un error en cursos, al hacer query, dando el error: %v", err)

	} else {
		defer rows.Close()
		for rows.Next() {
			var nombre string
			_ = rows.Scan(&nombre)

			opciones.Push(
				NewOpcion(nombre, fmt.Sprintf("/Cursos?path=%s", nombre)),
			)
		}

		return NewData(NewContenidoMinimo(c.PathActual(), returnPath), opciones.Items()), nil
	}
}

func (c *Cursos) PathActual() string {
	if elemento, err := c.Path.Desapilar(); err != nil {
		return "Cursos"

	} else {
		pathActual := fmt.Sprintf("%s > %s", c.PathActual(), elemento)
		c.Path.Apilar(elemento)
		return pathActual
	}
}

func (c *Cursos) Cd(subpath string) error {
	if subpath == "" {
		return nil
	}

	if subpath == ".." {
		return c.RutinaAtras()
	}

	query := ""
	switch c.Tipo {
	case TCC_GENERAL:
		query = fmt.Sprintf(QUERY_OBTENER_CURSO, subpath, subpath)
	case TCC_DENTRO_CURSO:
		if idCurso, err := c.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTENER_TEMA_CURSO, idCurso, subpath)
		}
	case TCC_DENTRO_TEMA:
		return fmt.Errorf("ya se esta viendo todos los archivos, no hay subcarpetas")
	}

	if query == "" {
		return fmt.Errorf("hubo un error en la query, y esta vacia")
	}

	fila := c.Bdd.QueryRow(query)
	var id int64
	var nombre string
	if err := fila.Scan(&id, &nombre); err != nil {
		return fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	switch c.Tipo {
	case TCC_GENERAL:
		c.Tipo = TCC_DENTRO_CURSO
	case TCC_DENTRO_CURSO:
		c.Tipo = TCC_DENTRO_TEMA
	}

	c.Indice.Apilar(id)
	c.Path.Apilar(nombre)
	return nil
}

func (c *Cursos) RutinaAtras() error {
	_, _ = c.Indice.Desapilar()
	_, _ = c.Path.Desapilar()

	switch c.Tipo {
	case TCC_GENERAL:
		return fmt.Errorf("no deberia ser posible que pongan .. aca")

	case TCC_DENTRO_CURSO:
		c.Tipo = TCC_GENERAL
	case TCC_DENTRO_TEMA:
		c.Tipo = TCC_DENTRO_CURSO
	}

	return nil
}
