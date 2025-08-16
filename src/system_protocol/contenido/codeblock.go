package contenido

// Lenguaje, codigo
type CodeBlock struct {
	NodoHeader
	Lenguaje string `bson:"lenguaje"`
	Codigo   string `bson:"codigo"`
}

func NewCodeBlock(lenguaje, codigo string) CodeBlock {
	return CodeBlock{
		NodoHeader: NewTipo(MK_CodeBlock),
		Lenguaje:   lenguaje,
		Codigo:     codigo,
	}
}
