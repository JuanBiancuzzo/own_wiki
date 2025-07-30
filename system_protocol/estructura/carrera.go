package estructura

import (
	"database/sql"
	"fmt"
)

const QUERY_CARRERA = "SELECT id FROM carreras WHERE nombre = ?"
const QUERY_CARRERA_PATH = `SELECT res.id FROM (
	SELECT carreras.id, archivos.path FROM archivos INNER JOIN carreras ON archivos.id = carreras.idArchivo
) AS res WHERE res.path = ?`
const INSERTAR_CARRERA = "INSERT INTO carreras (nombre, etapa, tieneCodigoMateria, idArchivo) VALUES (?, ?, ?, ?)"

type Carrera struct {
	PathArchivo string
	Nombre      string
	Etapa       Etapa
	TieneCodigo bool
}

func NewCarrera(pathArchivo string, nombre string, repEtapa string, tieneCodigo string) (*Carrera, error) {
	if etapa, err := ObtenerEtapa(repEtapa); err != nil {
		return nil, fmt.Errorf("error al crear carrera con error: %v", err)
	} else {
		return &Carrera{
			PathArchivo: pathArchivo,
			Nombre:      nombre,
			Etapa:       etapa,
			TieneCodigo: BooleanoODefault(tieneCodigo, false),
		}, nil
	}
}

func (c *Carrera) Insertar(idArchivo int64) []any {
	return []any{
		c.Nombre,
		c.Etapa,
		c.TieneCodigo,
		idArchivo,
	}
}

func (c *Carrera) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- fmt.Sprintf("Insertar Carrera: %s", c.Nombre)

	if idArchivo, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_ARCHIVO, c.PathArchivo) },
	); !existe {
		return false

	} else if _, err := bdd.Exec(INSERTAR_CARRERA, c.Insertar(idArchivo)...); err != nil {
		canal <- fmt.Sprintf("error al insertar una carrera, con error: %v", err)
	}

	return true
}
