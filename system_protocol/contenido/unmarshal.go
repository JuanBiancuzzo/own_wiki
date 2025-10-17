package contenido

import (
	"bytes"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Hijos) UnmarshalBSONValue(typ bson.Type, d []byte) error {
	if typ != bson.TypeArray {
		return fmt.Errorf("este tipo de dato no es un array")
	}

	rawArray, err := bson.ReadArray(bytes.NewReader(d))
	if err != nil {
		return err
	}

	valores, err := rawArray.Values()
	if err != nil {
		return err
	}

	hijos := make([]any, len(valores))
	for i, raw := range valores {
		var header NodoHeader
		if err := bson.Unmarshal(raw.Value, &header); err != nil {
			return err
		}

		if hijo, err := ObtenerEstructura(header.Tipo, raw.Value); err != nil {
			return err

		} else {
			hijos[i] = hijo
		}
	}

	return nil
}

/*
func (h *Hijos) UnmarshalJSON(d []byte) error {
	dec := json.NewDecoder(bytes.NewReader(d))
	if token, err := dec.Token(); err != nil {
		return err

	} else if delim, ok := token.(json.Delim); !ok {
		return fmt.Errorf("primer token (%v) no es un delimitador", token)

	} else if delim.String() != "[" {
		return fmt.Errorf("el delimitador inicial no es [ sino que %s", delim.String())
	}

	var header NodoHeader
	var raw json.RawMessage
	hijos := []any{}
	for dec.More() {
		if err := dec.Decode(&raw); err != nil {
			return err
		}
		if err := json.Unmarshal(raw, &header); err != nil {
			return err
		}

		if hijo, err := ObtenerEstructura(header.Tipo, raw); err != nil {
			return err

		} else {
			hijos = append(hijos, hijo)
		}
	}

	*h = hijos
	return nil
}
*/

func ObtenerEstructura(tipo MkTipo, d []byte) (any, error) {
	switch tipo {
	case MK_Header1:
		fallthrough
	case MK_Header2:
		fallthrough
	case MK_Header3:
		fallthrough
	case MK_Header4:
		fallthrough
	case MK_Header5:
		fallthrough
	case MK_Header6:
		var header Header
		return header, bson.Unmarshal(d, &header)

	case MK_CodeBlock:
		var codeBlock CodeBlock
		return codeBlock, bson.Unmarshal(d, &codeBlock)

	case MK_Image:
		var imagen Imagen
		return imagen, bson.Unmarshal(d, &imagen)

	case MK_Link:
		var link Link
		return link, bson.Unmarshal(d, &link)

	case MK_ListDesordenada:
		var lista ListaDesordenada
		return lista, bson.Unmarshal(d, &lista)

	case MK_ListOrdenada:
		var lista ListaOrdenada
		return lista, bson.Unmarshal(d, &lista)

	case MK_Table:
		var tabla Tabla
		return tabla, bson.Unmarshal(d, &tabla)

	case MK_Callout:
		var callout Callout
		return callout, bson.Unmarshal(d, &callout)

	case MK_Code:
		fallthrough
	case MK_HtmlBlock:
		fallthrough
	case MK_HtmlSpan:
		fallthrough
	case MK_Math:
		fallthrough
	case MK_MathBlock:
		fallthrough
	case MK_Subscript:
		fallthrough
	case MK_Superscript:
		fallthrough
	case MK_Text:
		var hoja HojaTexto
		return hoja, bson.Unmarshal(d, &hoja)

	case MK_Hardbreak:
		fallthrough
	case MK_HorizontalRule:
		var hoja HojaVacia
		return hoja, bson.Unmarshal(d, &hoja)

	case MK_Aside:
		fallthrough
	case MK_BlockQuote:
		fallthrough
	case MK_Parrafo:
		fallthrough
	case MK_Bold:
		fallthrough
	case MK_Italic:
		var contenedor Contenedor
		return contenedor, bson.Unmarshal(d, &contenedor)
	}

	return nil, fmt.Errorf("no existe la forma de transformar el tipo '%s'", tipo.String())
}
