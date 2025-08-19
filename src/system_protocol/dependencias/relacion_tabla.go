package dependencias

type RelacionTabla struct {
	Tabla string
	Datos ConjuntoDato
}

func NewRelacion(tabla string, datos ConjuntoDato) RelacionTabla {
	return RelacionTabla{
		Tabla: tabla,
		Datos: datos,
	}
}
