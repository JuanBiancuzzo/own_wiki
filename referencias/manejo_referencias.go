package referencias

import (
	sp "own_wiki/system_protocol"
)

type ManejoReferencias struct {
	Referencias map[uint64]Referencia
	Contador    *sp.ContadorGen[uint64]
}

func incrementarNumReferencia(numero uint64) uint64 {
	return numero + 1
}

func decrementarNumReferencia(numero uint64) uint64 {
	return numero - 1
}

func NewManejoReferencias() *ManejoReferencias {
	incrementar := sp.Inc[uint64](incrementarNumReferencia)
	decrementar := sp.Dec[uint64](decrementarNumReferencia)

	return &ManejoReferencias{
		Referencias: make(map[uint64](Referencia)),
		Contador:    sp.NewContadorGen(incrementar, decrementar, 0),
	}
}
