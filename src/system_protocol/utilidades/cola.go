package utilidades

type Cola[T any] struct {
	Lista *Lista[T]
}

func NewCola[T any]() *Cola[T] {
	return &Cola[T]{
		Lista: NewLista[T](),
	}
}

func NewColaConCapacidad[T any](capacidad uint32) *Cola[T] {
	return &Cola[T]{
		Lista: NewListaConCapacidad[T](capacidad),
	}
}

func (c *Cola[T]) Encolar(elemento T) {
	c.Lista.Push(elemento)
}

func (c *Cola[T]) Desencolar() (T, error) {
	return c.Lista.SacarEn(0)
}

func (c *Cola[T]) Vacia() bool {
	return c.Lista.Largo == 0
}
