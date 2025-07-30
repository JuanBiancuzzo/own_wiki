package estructura

import (
	"database/sql"
	"fmt"
)

const INSERTAR_MATERIA_EQUIVALENTES = "INSERT INTO materiasEquivalentes (nombre, codigo, idMateria, idArchivo) VALUES (?, ?, ?, ?)"
const QUERY_MATERIA_EQUIVALENTES_PATH = `SELECT res.id FROM (
	SELECT materiasEquivalentes.id, archivos.path FROM archivos INNER JOIN materiasEquivalentes ON archivos.id = materiasEquivalentes.idArchivo
) AS res WHERE res.path = ?`

type MateriaEquivalente struct {
	PathArchivo string
	PathMateria string
	Nombre      string
	Codigo      string
}

func NewMateriaEquivalente(pathArchivo string, pathMateria string, nombre string, codigo string) *MateriaEquivalente {
	return &MateriaEquivalente{
		PathArchivo: pathArchivo,
		PathMateria: pathMateria,
		Nombre:      nombre,
		Codigo:      codigo,
	}
}

func (me *MateriaEquivalente) Insertar(idMateria int64, idArchivo int64) []any {
	return []any{
		me.Nombre,
		me.Codigo,
		idMateria,
		idArchivo,
	}
}

func (me *MateriaEquivalente) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- fmt.Sprintf("Insertar Materia Correlativas: %s => %s", me.Nombre, me.PathMateria)

	if idArchivo, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_ARCHIVO, me.PathArchivo) },
	); !existe {
		return false

	} else if idMateria, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_MATERIA_PATH, me.PathMateria) },
	); !existe {
		return false

	} else if _, err := bdd.Exec(INSERTAR_MATERIA_CORRELATIVA, me.Insertar(idMateria, idArchivo)...); err != nil {
		canal <- fmt.Sprintf("error al insertar una materia equivalente, con error: %v", err)
	}

	return true
}
