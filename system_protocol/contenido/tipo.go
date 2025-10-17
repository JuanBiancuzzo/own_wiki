package contenido

type MkTipo byte

func (tipo MkTipo) String() string {
	switch tipo {
	case MK_Header1:
		return "Header 1"
	case MK_Header2:
		return "Header 2"
	case MK_Header3:
		return "Header 3"
	case MK_Header4:
		return "Header 4"
	case MK_Header5:
		return "Header 5"
	case MK_Header6:
		return "Header 6"
	case MK_CodeBlock:
		return "CodeBlock"
	case MK_Image:
		return "Image"
	case MK_Link:
		return "Link"
	case MK_ListOrdenada:
		return "List Ordenada"
	case MK_ListDesordenada:
		return "List Desordenada"
	case MK_Callout:
		return "Callout"
	case MK_Code:
		return "Code"
	case MK_HtmlBlock:
		return "HtmlBlock"
	case MK_HtmlSpan:
		return "HtmlSpan"
	case MK_Hardbreak:
		return "Hardbreak"
	case MK_HorizontalRule:
		return "HorizontalRule"
	case MK_Math:
		return "Math"
	case MK_Subscript:
		return "Subscript"
	case MK_Superscript:
		return "Superscript"
	case MK_Text:
		return "Text"
	case MK_Aside:
		return "Aside"
	case MK_BlockQuote:
		return "BlockQuote"
	case MK_MathBlock:
		return "MathBlock"
	case MK_Parrafo:
		return "Parrafo"
	case MK_Bold:
		return "Bold"
	case MK_Italic:
		return "Italic"
	}

	return "[[ERROR]]"
}

const (
	// Contenedor y leaf
	MK_Header1 = iota
	MK_Header2
	MK_Header3
	MK_Header4
	MK_Header5
	MK_Header6
	MK_CodeBlock
	MK_Image
	MK_Link
	MK_ListOrdenada
	MK_ListDesordenada
	MK_Table

	MK_Callout
	MK_Wikilink

	// Leaf
	MK_Code           // solo texto
	MK_HtmlBlock      // Solo texto
	MK_HtmlSpan       // Solo texto
	MK_Hardbreak      // nada
	MK_HorizontalRule // Nada
	MK_Math           // Solo texto
	MK_Subscript      // Solo texto
	MK_Superscript    // Solo texto
	MK_Text           // solo texto

	// Contenedores:
	MK_Aside
	MK_BlockQuote
	MK_MathBlock
	MK_Parrafo
	MK_Bold
	MK_Italic
)

type NodoHeader struct {
	Tipo MkTipo `bson:"tipo"`
}

func NewTipo(tipo MkTipo) NodoHeader {
	return NodoHeader{
		Tipo: tipo,
	}
}
