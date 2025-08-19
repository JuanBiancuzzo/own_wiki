package views

import (
	"fmt"
)

type PedidoElementos struct {
	Endpoint string
	Cantidad int
	Offset   int
}

func CreateURLPedido(pedido *PedidoElementos, cantidad int) string {
	if cantidad <= 0 {
		cantidad = pedido.Cantidad
	}

	return fmt.Sprintf("/%s?cantidad=%d&offset=%d", pedido.Endpoint, cantidad, pedido.Offset)
}
