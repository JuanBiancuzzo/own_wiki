package dependencias

type RelacionTabla struct {
	Tabla           string
	Clave           string
	InfoRelacionada []RelacionTabla
	Datos           []any
}

func NewRelacionSimple(tabla, clave string, datos ...any) RelacionTabla {
	return RelacionTabla{
		Tabla:           tabla,
		Clave:           clave,
		InfoRelacionada: []RelacionTabla{},
		Datos:           datos,
	}
}

func NewRelacionCompleja(tabla, clave string, infoRelaciones []RelacionTabla, datos ...any) RelacionTabla {
	return RelacionTabla{
		Tabla:           tabla,
		Clave:           clave,
		InfoRelacionada: infoRelaciones,
		Datos:           datos,
	}
}
