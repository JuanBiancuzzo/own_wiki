package views

import (
	"fmt"
)

const CANTIDAD_DEFAULT = 5

type PedidoElementos struct {
	Endpoint string
	Cantidad int
	Offset   int
}

func NewPedidoElementos(cantidad int) PedidoElementos {
	if cantidad <= 0 {
		cantidad = CANTIDAD_DEFAULT
	}

	return PedidoElementos{
		Cantidad: cantidad,
		Offset:   0,
	}
}

func CreateURLPedido(pedido *PedidoElementos, cantidad int) string {
	if cantidad <= 0 {
		cantidad = pedido.Cantidad
	}

	return fmt.Sprintf("/%s?cantidad=%d&offset=%d", pedido.Endpoint, cantidad, pedido.Offset)
}
