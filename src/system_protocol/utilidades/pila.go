package utilidades

import "fmt"

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

func (p *Pila[T]) Apilar(elemento T) {
	p.Lista.Push(elemento)
}

func (p *Pila[T]) Pick() (T, error) {
	if p.Lista.Largo == 0 {
		return valor_default[T](), fmt.Errorf("no hay elemento para ver")
	}

	return p.Lista.ObtenerEn(p.Lista.Largo - 1)
}

func (p *Pila[T]) Desapilar() (T, error) {
	return p.Lista.Pop()
}

func (p *Pila[T]) Vacia() bool {
	return p.Lista.Largo == 0
}

func (p *Pila[T]) DesapilarIterativamente(yield func(T) bool) {
	for !p.Vacia() {
		if elemento, err := p.Desapilar(); err != nil {
			return
		} else if !yield(elemento) {
			return
		}
	}
}
