package dependencias

type ReferenciaTabla struct {
	Representativo bool
	Clave          string
	Tablas         []*DescripcionTabla
}

func NewReferenciaTabla(clave string, tablas []*DescripcionTabla, representativo bool) ReferenciaTabla {
	return ReferenciaTabla{
		Representativo: representativo,
		Clave:          clave,
		Tablas:         tablas,
	}
}
