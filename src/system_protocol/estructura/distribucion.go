package estructura

import (
	"database/sql"
	"fmt"
)

type TipoDistribucion string

const (
	DISTRIBUCION_DISCRETA     = "Discreta"
	DISTRIBUCION_CONTINUA     = "Continua"
	DISTRIBUCION_MULTIVARIADA = "Multivariada"
)

const INSERTAR_DISTRIBUCION = "INSERT INTO distribuciones (nombre, tipo, idArchivo) VALUES (?, ?, ?)"

type Distribucion struct {
	Nombre    string
	Tipo      TipoDistribucion
	IdArchivo *Opcional[int64]
}

func NewDistribucion(nombre string, repTipo string) (*Distribucion, error) {
	if tipo, err := ObtenerTipoDistribucion(repTipo); err != nil {
		return nil, fmt.Errorf("error al crear distribucion con error: %v", err)

	} else {
		return &Distribucion{
			Tipo:      tipo,
			Nombre:    nombre,
			IdArchivo: NewOpcional[int64](),
		}, nil
	}
}

func (c *Distribucion) CrearDependenciaArchivo(dependible Dependible) {
	dependible.CargarDependencia(func(id int64) (Cargable, bool) {
		c.IdArchivo.Asignar(id)
		return c, true
	})
}

func (c *Distribucion) Insertar() ([]any, error) {
	if idArchivo, existe := c.IdArchivo.Obtener(); !existe {
		return []any{}, fmt.Errorf("distribucion no tiene todavia el idArchivo")
	} else {
		return []any{c.Nombre, c.Tipo, idArchivo}, nil
	}
}

func (c *Distribucion) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	canal <- "Insertar Distribucion"
	if datos, err := c.Insertar(); err != nil {
		return 0, err
	} else {
		return InsertarDirecto(bdd, INSERTAR_DISTRIBUCION, datos...)
	}
}

func (c *Distribucion) ResolverDependencias(id int64) []Cargable {
	return []Cargable{}
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
