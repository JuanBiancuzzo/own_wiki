package fs

import (
	"database/sql"
	"fmt"
	"strings"

	u "own_wiki/system_protocol/utilidades"

	"github.com/labstack/echo/v4"
)

const (
	QUERY_CARRERAS_LS = "SELECT nombre FROM carreras"
	QUERY_MATERIAS_LS = `SELECT materiasGlobal.nombre FROM (
		SELECT idCarrera, nombre FROM materias UNION ALL SELECT idCarrera, nombre FROM materiasEquivalentes
	) AS materiasGlobal WHERE materiasGlobal.idCarrera = %d`
	QUERY_TEMAS_MATERIA_LS = "SELECT nombre FROM temasMateria WHERE idMateria = %d ORDER BY capitulo, parte"
	QUERY_NOTA_MATERIA_LS  = `SELECT DISTINCT notas.nombre FROM notas INNER JOIN (
		SELECT idNota, idVinculo FROM notasVinculo WHERE tipoVinculo = "Facultad"
	) AS vinculo ON notas.id = vinculo.idNota WHERE vinculo.idVinculo = %d`
)

const (
	QUERY_OBTENER_CARRERA = "SELECT id, nombre FROM carreras WHERE nombre = '%s'"
	QUERY_OBTENER_MATERIA = `SELECT id, nombre FROM materias WHERE idCarrera = %d AND nombre = '%s'
		UNION ALL
	SELECT idMateria AS id, nombre FROM materiasEquivalentes WHERE idCarrera = %d AND nombre = '%s'
	`
	QUERY_OBTNER_TEMA_MATERIA = "SELECT id, nombre FROM temasMateria WHERE idMateria = %d AND nombre = '%s'"
)

type TipoFacultad byte

const (
	TF_FACULTAD = iota
	TF_DENTRO_CARRERA
	TF_DENTRO_MATERIA
	TF_DENTRO_TEMA
)

type Facultad struct {
	Bdd    *sql.DB
	Tipo   TipoFacultad
	Indice *u.Pila[int64]
	Path   *u.Pila[string]
}

func NewFacultad(bdd *sql.DB) *Facultad {
	return &Facultad{
		Bdd:    bdd,
		Tipo:   TF_FACULTAD,
		Indice: u.NewPila[int64](),
		Path:   u.NewPila[string](),
	}
}

func (f *Facultad) DeterminarRuta(ec echo.Context) error {
	path := strings.TrimSpace(ec.QueryParam("path"))

	var carpetaActual string
	var errCd error
	for subpath := range strings.SplitSeq(path, "/") {
		carpetaActual, errCd = f.Cd(subpath)
		if errCd != nil {
			break
		}
		if carpetaActual == PD_ROOT {
			return ec.Render(200, "root", DATA_ROOT)
		}
	}

	if carpetaActual == "" {
		carpetaActual = "Facultad"
	}
	data, errLs := f.Ls(carpetaActual)

	if errCd != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Cd tuvo el error: %v", errCd), "/Facultad"))
	}
	if errLs != nil {
		data.Opciones = append(data.Opciones, NewOpcion(fmt.Sprintf("Ls tuvo el error: %v", errLs), "/Facultad"))
	}

	return ec.Render(200, "facultad", data)
}

func (f *Facultad) Ls(carpetaActual string) (Data, error) {
	var data Data
	opciones := u.NewLista[Opcion]()
	returnPath := "/Facultad?path=.."

	var query string
	switch f.Tipo {
	case TF_FACULTAD:
		query = QUERY_CARRERAS_LS
		returnPath = "/Root"
	case TF_DENTRO_CARRERA:
		if idCarrera, err := f.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_MATERIAS_LS, idCarrera)
		}
	case TF_DENTRO_MATERIA:
		if idMateria, err := f.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_TEMAS_MATERIA_LS, idMateria)
		}
	case TF_DENTRO_TEMA:
		if idTema, err := f.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_NOTA_MATERIA_LS, idTema)
		}

	}

	if rows, err := f.Bdd.Query(query); err != nil {
		return data, fmt.Errorf("se obtuvo un error en facultad, al hacer query, dando el error: %v", err)

	} else {
		defer rows.Close()
		for rows.Next() {
			var nombre string
			_ = rows.Scan(&nombre)

			opciones.Push(
				NewOpcion(nombre, fmt.Sprintf("/Facultad?path=%s", nombre)),
			)
		}

		return NewData(
			NewTextoVinculo(carpetaActual, returnPath),
			f.PathActual(0),
			opciones.Items(),
		), nil
	}
}

func (f *Facultad) PathActual(profundidad int) []TextoVinculo {
	if profundidad > 2 {
		return []TextoVinculo{
			NewTextoVinculo("...", fmt.Sprintf("/Facultad?path=%s", strings.Repeat("../", profundidad))),
		}
	}

	if elemento, err := f.Path.Desapilar(); err != nil {
		return []TextoVinculo{
			NewTextoVinculo("Own_wiki", fmt.Sprintf("/Facultad?path=%s", strings.Repeat("../", profundidad+1))),
			NewTextoVinculo("Facultad", fmt.Sprintf("/Facultad?path=%s", strings.Repeat("../", profundidad))),
		}
	} else {
		textoVinculo := NewTextoVinculo(elemento, fmt.Sprintf("/Facultad?path=%s", strings.Repeat("../", profundidad)))
		pathActual := append(f.PathActual(profundidad+1), textoVinculo)
		f.Path.Apilar(elemento)
		return pathActual
	}
}

func (f *Facultad) Cd(subpath string) (string, error) {
	if subpath == "" {
		carpetaActual, _ := f.Path.Pick()
		return carpetaActual, nil
	}

	if subpath == ".." {
		return f.RutinaAtras()
	}

	query := ""
	switch f.Tipo {
	case TF_FACULTAD:
		query = fmt.Sprintf(QUERY_OBTENER_CARRERA, subpath)
	case TF_DENTRO_CARRERA:
		if idCarrera, err := f.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTENER_MATERIA, idCarrera, subpath, idCarrera, subpath)
		}
	case TF_DENTRO_MATERIA:
		if idMateria, err := f.Indice.Pick(); err == nil {
			query = fmt.Sprintf(QUERY_OBTNER_TEMA_MATERIA, idMateria, subpath)
		}
	case TF_DENTRO_TEMA:
		return "", fmt.Errorf("ya se esta viendo todos los archivos, no hay subcarpetas")
	}

	if query == "" {
		return "", fmt.Errorf("hubo un error en la query, y esta vacia")
	}

	fila := f.Bdd.QueryRow(query)
	var id int64
	var nombre string

	if err := fila.Scan(&id, &nombre); err != nil {
		return "", fmt.Errorf("no existe posible solucion para el cd a '%s', con error: %v", subpath, err)
	}

	switch f.Tipo {
	case TF_FACULTAD:
		f.Tipo = TF_DENTRO_CARRERA
	case TF_DENTRO_CARRERA:
		f.Tipo = TF_DENTRO_MATERIA
	case TF_DENTRO_MATERIA:
		f.Tipo = TF_DENTRO_TEMA
	}

	f.Indice.Apilar(id)
	f.Path.Apilar(nombre)
	return subpath, nil
}

func (f *Facultad) RutinaAtras() (string, error) {
	_, _ = f.Indice.Desapilar()

	switch f.Tipo {
	case TF_FACULTAD:
		return PD_ROOT, nil

	case TF_DENTRO_CARRERA:
		f.Tipo = TF_FACULTAD
	case TF_DENTRO_MATERIA:
		f.Tipo = TF_DENTRO_CARRERA
	case TF_DENTRO_TEMA:
		f.Tipo = TF_DENTRO_MATERIA
	}

	return f.Path.Desapilar()
}
