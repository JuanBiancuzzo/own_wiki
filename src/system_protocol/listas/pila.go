package listas

type Pila[T any] struct {
	Lista *Lista[T]
}

func NewPila[T any]() *Pila[T] {
	return &Pila[T]{
		Lista: NewLista[T](),
	}
}

func NewPilaConCapacidad[T any](capacidad uint32) *Pila[T] {
	return &Pila[T]{
		Lista: NewListaConCapacidad[T](capacidad),
	}
}

func (l *Pila[T]) Apilar(elemento T) {
	l.Lista.Push(elemento)
}

func (l *Pila[T]) Desapilar() (T, error) {
	return l.Lista.Pop()
}

func (l *Pila[T]) Vacia() bool {
	return l.Lista.Largo == 0
}
