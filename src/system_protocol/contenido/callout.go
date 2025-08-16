package contenido

// Tipo de callout (el ![info] seria info), titulo, y texto
type Callout struct {
	NodoHeader
	Tipo      string `bson:"tipoCallout"` // El ![tipo]
	Titulo    string `bson:"titulo"`
	Elementos Hijos  `bson:"elementos"`
}

func NewCallout(tipo, titulo string, elementos Hijos) Callout {
	return Callout{
		NodoHeader: NewTipo(MK_Callout),
		Tipo:       tipo,
		Titulo:     titulo,
		Elementos:  elementos,
	}
}
