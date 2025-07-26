package referencias

import (
	sp "own_wiki/system_protocol"
)

type Referencia interface {
	Modificar(contadorNumReferencia *sp.ContadorGen[uint64])
}
