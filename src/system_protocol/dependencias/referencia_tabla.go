package dependencias

type ReferenciaTabla struct {
	Representativo bool
	Clave          string
	Tabla          DescripcionTabla
}

func NewReferenciaTabla(clave string, tabla DescripcionTabla, representativo bool) ReferenciaTabla {
	return ReferenciaTabla{
		Representativo: representativo,
		Clave:          clave,
		Tabla:          tabla,
	}
}
