package system_protocol

import (
	"cmp"
	"fmt"

	l "own_wiki/system_protocol/listas"
)

type Inc[T cmp.Ordered] func(T) T
type Dec[T cmp.Ordered] func(T) T

type ContadorGen[T cmp.Ordered] struct {
	Incrementar  Inc[T]
	Decrementar  Dec[T]
	UltimoNumero T
	Libres       *l.Lista[T] // Podriamos usar un heap
}

func NewContadorGen[T cmp.Ordered](incrementar Inc[T], decrementar Dec[T], inicio T) *ContadorGen[T] {
	return &ContadorGen[T]{
		Incrementar:  incrementar,
		Decrementar:  decrementar,
		UltimoNumero: inicio,
		Libres:       l.NewLista[T](),
	}
}

func (c *ContadorGen[T]) Obtener() T {
	elementoMasChico, err := c.Libres.SacarEn(0)
	if err != nil {
		elementoMasChico = c.UltimoNumero
		c.UltimoNumero = c.Incrementar(c.UltimoNumero)
	}

	return elementoMasChico
}

func (c *ContadorGen[T]) Devolver(numero T) error {
	if c.UltimoNumero == numero {
		c.UltimoNumero = c.Decrementar(c.UltimoNumero)
		return nil
	}

	if c.UltimoNumero < numero {
		return fmt.Errorf("El numero a devolver (%v) es mayor al numero mas grande guardado (%v)", numero, c.UltimoNumero)
	}

	// Insertamos de forma ordenada para hacer el obtener mucho mas rapido
	indiceInsertar := 0
	for i, numeroReservador := range c.Libres.Items() {
		if numeroReservador < numero {
			continue
		}

		if numeroReservador == numero {
			return fmt.Errorf("Se esta devolviendo (%v) un numero que ya se devolvio", numero)
		}

		indiceInsertar = i
		break
	}

	c.Libres.AgregarEn(numero, uint32(indiceInsertar))
	return nil
}
