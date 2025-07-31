package estructura

import (
	"database/sql"
	"fmt"
)

const INSERTAR_MATERIA_CORRELATIVA = "INSERT INTO materiasCorrelativas (tipoMateria, idMateria, tipoCorrelativa, idCorrelativa) VALUES (?, ?, ?, ?)"

type TipoMateria string

const (
	MATERIA_REAL        = "Materia"
	MATERIA_EQUIVALENTE = "Equivalente"
)

type MateriasCorrelativas struct {
	PathArchivo     string
	TipoArchivo     TipoMateria
	PathCorrelativa string
	TipoCorrelativa TipoMateria
}

func NewMateriasCorrelativas(pathArchivo string, tipoArchivo TipoMateria, pathCorrelativa string, tipoCorrelativa TipoMateria) *MateriasCorrelativas {
	return &MateriasCorrelativas{
		PathArchivo:     pathArchivo,
		TipoArchivo:     tipoArchivo,
		PathCorrelativa: pathCorrelativa,
		TipoCorrelativa: tipoCorrelativa,
	}
}

func (mc *MateriasCorrelativas) Insertar(idMateria int64, idCorrelativa int64) []any {
	return []any{
		mc.TipoArchivo,
		idMateria,
		mc.TipoCorrelativa,
		idCorrelativa,
	}
}

func (mc *MateriasCorrelativas) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- fmt.Sprintf("Insertar Materia Correlativas entre: %s(%s) => %s(%s)", Nombre(mc.PathArchivo), mc.TipoArchivo, Nombre(mc.PathCorrelativa), mc.TipoCorrelativa)

	queryMateria := QUERY_MATERIA_PATH
	if mc.TipoArchivo == MATERIA_EQUIVALENTE {
		queryMateria = QUERY_MATERIA_EQUIVALENTES_PATH
	}
	queryCorrelativa := QUERY_MATERIA_PATH
	if mc.TipoCorrelativa == MATERIA_EQUIVALENTE {
		queryCorrelativa = QUERY_MATERIA_EQUIVALENTES_PATH
	}
	if idMateria, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(queryMateria, mc.PathArchivo) },
	); !existe {
		return false

	} else if idCorrelativa, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(queryCorrelativa, mc.PathCorrelativa) },
	); !existe {
		return false

	} else if _, err := bdd.Exec(INSERTAR_MATERIA_CORRELATIVA, mc.Insertar(idMateria, idCorrelativa)...); err != nil {
		canal <- fmt.Sprintf("error al insertar una materias correlativas, con error: %v", err)
	}

	return true
}
