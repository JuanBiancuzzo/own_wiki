package contenido

/*
	MK_Aside
	MK_BlockQuote
	MK_Parrafo
	MK_Bold
	MK_Italic
*/
type Contenedor struct {
	NodoHeader
	Elementos Hijos `bson:"elementos"`
}

func newContenedor(tipo MkTipo, elementos Hijos) Contenedor {
	return Contenedor{
		NodoHeader: NewTipo(tipo),
		Elementos:  elementos,
	}
}

func NewAside(elementos Hijos) Contenedor      { return newContenedor(MK_Aside, elementos) }
func NewBlockQuote(elementos Hijos) Contenedor { return newContenedor(MK_BlockQuote, elementos) }
func NewParrafo(elementos Hijos) Contenedor    { return newContenedor(MK_Parrafo, elementos) }
func NewBold(elementos Hijos) Contenedor       { return newContenedor(MK_Bold, elementos) }
func NewItalic(elementos Hijos) Contenedor     { return newContenedor(MK_Italic, elementos) }
