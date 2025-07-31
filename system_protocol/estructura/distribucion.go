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

type ConstructorDistribucion struct {
	Tipo   TipoDistribucion
	Nombre string
}

func NewConstructorDistribucion(nombre string, repTipo string) (*ConstructorDistribucion, error) {
	if tipo, err := ObtenerTipoDistribucion(repTipo); err != nil {
		return nil, fmt.Errorf("error al crear distribucion con error: %v", err)

	} else {
		return &ConstructorDistribucion{
			Tipo:   tipo,
			Nombre: nombre,
		}, nil
	}
}

func (cd *ConstructorDistribucion) CumpleDependencia(id int64) (Cargable, bool) {
	return &Distribucion{
		Nombre:    cd.Nombre,
		Tipo:      cd.Tipo,
		IdArchivo: id,
	}, true
}

type Distribucion struct {
	Nombre    string
	Tipo      TipoDistribucion
	IdArchivo int64
}

func (d *Distribucion) Insertar() []any {
	return []any{
		d.Nombre,
		d.Tipo,
		d.IdArchivo,
	}
}

func (d *Distribucion) CargarDatos(bdd *sql.DB, canal chan string) (int64, error) {
	// canal <- "Insertar Distribucion"
	return Insertar(
		func() (sql.Result, error) { return bdd.Exec(INSERTAR_DISTRIBUCION, d.Insertar()...) },
	)
}

func (d *Distribucion) ResolverDependencias(id int64) []Cargable {
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
