package contenido

/*
	MK_Code           // solo texto
	MK_HtmlBlock      // Solo texto
	MK_HtmlSpan       // Solo texto
	MK_Math           // Solo texto
	MK_MathBlock // Solo texto
	MK_Subscript      // Solo texto
	MK_Superscript    // Solo texto
	MK_Text           // solo texto
*/
type HojaTexto struct {
	NodoHeader
	Texto string `bson:"texto"`
}

func newHojaTexto(tipo MkTipo, texto string) HojaTexto {
	return HojaTexto{
		NodoHeader: NewTipo(tipo),
		Texto:      texto,
	}
}

func NewCode(texto string) HojaTexto        { return newHojaTexto(MK_Code, texto) }
func NewHtmlBlock(texto string) HojaTexto   { return newHojaTexto(MK_HtmlBlock, texto) }
func NewHtmlSpan(texto string) HojaTexto    { return newHojaTexto(MK_HtmlSpan, texto) }
func NewMath(texto string) HojaTexto        { return newHojaTexto(MK_Math, texto) }
func NewMathBlock(texto string) HojaTexto   { return newHojaTexto(MK_MathBlock, texto) }
func NewSubscript(texto string) HojaTexto   { return newHojaTexto(MK_Subscript, texto) }
func NewSuperscript(texto string) HojaTexto { return newHojaTexto(MK_Superscript, texto) }
func NewText(texto string) HojaTexto        { return newHojaTexto(MK_Text, texto) }

/*
	MK_Hardbreak      // nada
	MK_HorizontalRule // Nada
*/
type HojaVacia struct {
	NodoHeader
}

func newHojaVacia(tipo MkTipo) HojaVacia {
	return HojaVacia{
		NodoHeader: NewTipo(tipo),
	}
}

func NewHardbreak() HojaVacia      { return newHojaVacia(MK_Hardbreak) }
func NewHorizontalRule() HojaVacia { return newHojaVacia(MK_HorizontalRule) }
