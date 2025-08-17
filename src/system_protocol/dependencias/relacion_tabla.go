package dependencias

type RelacionTabla struct {
	Tabla           string
	InfoRelacionada []RelacionTabla
	Datos           ConjuntoDato
}

func NewRelacion(tabla string, datos ConjuntoDato) RelacionTabla {
	return RelacionTabla{
		Tabla:           tabla,
		InfoRelacionada: []RelacionTabla{},
		Datos:           datos,
	}
}
