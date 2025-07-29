package estructura

import (
	"database/sql"
	"fmt"
)

const INSERTAR_DISTRIBUCION = "INSERT INTO distribuciones (nombre, tipo, idArchivo) VALUES (?, ?, ?)"

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

type Distribucion struct {
	PathArchivo string
	Tipo        TipoDistribucion
	Nombre      string
}

func NewDistribucion(pathArchivo string, nombre string, repTipo string) (*Distribucion, error) {
	if tipo, err := ObtenerTipoDistribucion(repTipo); err != nil {
		return nil, fmt.Errorf("error al crear distribucion con error: %v", err)

	} else {
		return &Distribucion{
			PathArchivo: pathArchivo,
			Tipo:        tipo,
			Nombre:      nombre,
		}, nil
	}
}

func ObtenerTipoDistribucion(representacion string) (TipoDistribucion, error) {
	var tipoDistribucion TipoDistribucion
	switch representacion {
	case "discreta":
		tipoDistribucion = DISTRIBUCION_DISCRETA
	case "continua":
		tipoDistribucion = DISTRIBUCION_CONTINUA
	case "multivariada":
		tipoDistribucion = DISTRIBUCION_MULTIVARIADA
	default:
		return DISTRIBUCION_DISCRETA, fmt.Errorf("el tipo de distribucion (%s) no es uno de los esperados", representacion)
	}

	return tipoDistribucion, nil
}

func (d *Distribucion) Insertar(idArchivo int64) []any {
	return []any{
		d.Nombre,
		d.Tipo,
		idArchivo,
	}
}

func (d *Distribucion) CargarDatos(bdd *sql.DB, canal chan string) bool {
	canal <- "Insertar Distribucion"

	if idArchivo, existe := Obtener(
		func() *sql.Row { return bdd.QueryRow(QUERY_ARCHIVO, d.PathArchivo) },
	); !existe {
		return false

	} else if _, err := bdd.Exec(INSERTAR_DISTRIBUCION, d.Insertar(idArchivo)...); err != nil {
		canal <- fmt.Sprintf("error al insertar una distribucion, con error: %v", err)
	}

	return true
}
