package listas

import (
	"fmt"
)

type Lista[T any] struct {
	Elementos []T
	Capacidad uint32
	Largo     uint32
}

const DEFAULT_CAPACITY = 2
const CAPACITY_MULTIPLICATION = 2

func NewLista[T any]() *Lista[T] {
	return NewListaConCapacidad[T](DEFAULT_CAPACITY)
}

func NewListaConCapacidad[T any](capacidad uint32) *Lista[T] {
	elementos := make([]T, capacidad)
	for i := 0; i < int(capacidad); i++ {
		elementos[i] = valor_default[T]()
	}

	return &Lista[T]{
		Elementos: elementos,
		Capacidad: capacidad,
		Largo:     0,
	}
}

func valor_default[T any]() T {
	var result T
	return result
}

func (l *Lista[T]) expandir() {
	newCapacity := l.Capacidad * CAPACITY_MULTIPLICATION
	nuevosElementos := make([]T, newCapacity)

	for i := 0; i < int(l.Largo); i++ {
		nuevosElementos[i] = l.Elementos[i]
	}

	l.Elementos = nuevosElementos
	l.Capacidad = newCapacity
}

func (l *Lista[T]) Push(elemento T) {
	if l.Largo == l.Capacidad {
		l.expandir()
	}

	l.Elementos[l.Largo] = elemento
	l.Largo++
}

func (l *Lista[T]) AgregarEn(elemento T, en uint32) error {
	if l.Largo < en {
		return fmt.Errorf("el indice de %d es mayor al largo de la lista", en)
	}

	if l.Largo == l.Capacidad {
		l.expandir()
	}

	elementoAReemplazar := elemento
	for i := int(en); i <= int(l.Largo); i++ {
		temp := l.Elementos[i]
		l.Elementos[i] = elementoAReemplazar
		elementoAReemplazar = temp
	}

	l.Largo++
	return nil
}

func (l *Lista[T]) ObtenerEn(en uint32) (T, error) {
	if l.Largo <= en {
		return valor_default[T](), fmt.Errorf("el largo de la lista (%d) es menor a la posicion pasada: %d", l.Largo, en)
	}

	return l.Elementos[en], nil
}

func (l *Lista[T]) Pop() (T, error) {
	if l.Largo == 0 {
		return valor_default[T](), fmt.Errorf("la lista esta vacia")
	}

	ultimo := l.Elementos[l.Largo-1]
	l.Largo--
	return ultimo, nil
}

func (l *Lista[T]) SacarEn(en uint32) (T, error) {
	if l.Largo <= en {
		return valor_default[T](), fmt.Errorf("el largo de la lista (%d) es menor a la posicion pasada: %d", l.Largo, en)
	}

	resultado := l.Elementos[en]
	for i := int(en); i < int(l.Largo)-1; i++ {
		l.Elementos[i] = l.Elementos[i+1]
	}
	l.Largo--

	return resultado, nil

}

func (l *Lista[T]) Vaciar() {
	for i := 0; i < int(l.Largo); i++ {
		l.Elementos[i] = valor_default[T]()
	}
	l.Largo = 0
}

func (l *Lista[T]) Items() []T {
	slice := make([]T, l.Largo)
	for i := 0; i < int(l.Largo); i++ {
		slice[i] = l.Elementos[i]
	}
	return slice
}

func (l *Lista[T]) Vacia() bool {
	return l.Largo == 0
}
