package contenido

type MkTipo string

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
	MK_Header1         = "Header 1"
	MK_Header2         = "Header 2"
	MK_Header3         = "Header 3"
	MK_Header4         = "Header 4"
	MK_Header5         = "Header 5"
	MK_Header6         = "Header 6"
	MK_CodeBlock       = "CodeBlock"
	MK_Image           = "Image"
	MK_Link            = "Link"
	MK_ListDesordenada = "ListDesordenada"
	MK_ListOrdenada    = "ListOrdenada"
	MK_Table           = "Table"

	MK_Callout  = "Callout"
	MK_WikiLink = "Wiki link"

	// Leaf
	MK_Code           = "Code"             // solo texto
	MK_HtmlBlock      = "Html block"       // Solo texto
	MK_HtmlSpan       = "Html span"        // Solo texto
	MK_Math           = "Math"             // Solo texto
	MK_Subscript      = "Subscript"        // Solo texto
	MK_Superscript    = "Superscript"      // Solo texto
	MK_Text           = "Texto"            // solo texto
	MK_Hardbreak      = "Hardbrak"         // nada
	MK_HorizontalRule = "Linea horizontal" // Nada

	// Contenedores:
	MK_Aside      = "Aside"
	MK_BlockQuote = "BlockQuote"
	MK_MathBlock  = "Mathblock"
	MK_Parrafo    = "Parrafo"
	MK_Bold       = "Bold"
	MK_Italic     = "Italics"
)

/*
type MkTipo byte

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
*/

type NodoHeader struct {
	Tipo MkTipo `bson:"tipo"`
}

func NewTipo(tipo MkTipo) NodoHeader {
	return NodoHeader{
		Tipo: tipo,
	}
}
