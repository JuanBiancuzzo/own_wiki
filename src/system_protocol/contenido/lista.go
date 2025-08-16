package contenido

type TipoLista byte

const (
	TL_Alfabetico = iota
	TL_Numerico
)

// El tipo de lista, numero de inicio si es ordenada, elementos hijos (otras listas o items)
type ListaDesordenada struct {
	NodoHeader
	Items []any `bson:"items"`
}

type ListaOrdenada struct {
	NodoHeader
	Tipo   TipoLista `bson:"tipoLista"`
	Items  []any     `bson:"items"`
	Inicio int       `bson:"inicio"`
}

type Item struct {
	Elementos Hijos `bson:"elementos"`
}

func NewListaDesordenada(items []any) ListaDesordenada {
	return ListaDesordenada{
		NodoHeader: NewTipo(MK_ListOrdenada),
		Items:      items,
	}
}

func NewListaOrdenada(tipo TipoLista, items []any, inicio int) ListaOrdenada {
	return ListaOrdenada{
		NodoHeader: NewTipo(MK_ListDesordenada),
		Tipo:       tipo,
		Items:      items,
		Inicio:     inicio,
	}
}

func NewItem(elementos Hijos) any {
	if len(elementos) > 1 {
		return Item{
			Elementos: elementos,
		}
	} else {
		return elementos[0]
	}
}
