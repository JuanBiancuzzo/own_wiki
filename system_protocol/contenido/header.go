package contenido

import "fmt"

type Header struct {
	NodoHeader
	Elementos Hijos `json:"elementos"`
}

func NewHeader(nivel uint8, elementos Hijos) (Header, error) {
	header := Header{
		Elementos: elementos,
	}

	switch nivel {
	case 1:
		header.NodoHeader = NewTipo(MK_Header1)
	case 2:
		header.NodoHeader = NewTipo(MK_Header2)
	case 3:
		header.NodoHeader = NewTipo(MK_Header3)
	case 4:
		header.NodoHeader = NewTipo(MK_Header4)
	case 5:
		header.NodoHeader = NewTipo(MK_Header5)
	case 6:
		header.NodoHeader = NewTipo(MK_Header6)

	case 0:
		return header, fmt.Errorf("el nivel 0 no es un header valido")
	default:
		return header, fmt.Errorf("el nivel %d supera el nivel inferior de 6", nivel)
	}

	return header, nil
}
