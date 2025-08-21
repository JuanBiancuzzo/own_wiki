package configuracion

import "encoding/json"

type TipoParametro string

const (
	TP_ELEMENTOS_PARCIAL  = "ElementosParcial"
	TP_ELEMENTOS_COMPLETO = "ElementosCompleto"
	TP_ELEMENTO_UNICO     = "ElementoUnico"
)

type Informacion struct {
	Detalles any
}

type HeaderParametro struct {
	Tipo TipoParametro `json:"tipo"`
}

type ParametroElementoUnico struct {
	HeaderParametro
	Tabla string `json:"tabla"`
	Id    string `json:"id"`
}

type ParametroElementosCompleto struct {
	HeaderParametro
	Tablas  map[string]InformacionTabla `json:"elementos"`
	GroupBy []string                    `json:"groupBy"`
}

type InformacionTabla struct {
	Condicion string   `json:"condicion"`
	OrderBy   []string `json:"orderBy"`
}

type ParametroElementosParcial struct {
	ParametroElementosCompleto
	Elementos InformacionParcial `json:"elementos"`
}

type InformacionParcial struct {
	Nombre   string `json:"nombrePedido"`
	Bloque   string `json:"bloquesElementos"`
	Cantidad int    `json:"cantidad"`
}

func (p *Informacion) UnmarshalJSON(d []byte) error {
	var header HeaderParametro
	if err := json.Unmarshal(d, &header); err != nil {
		return err
	}

	switch header.Tipo {
	case TP_ELEMENTO_UNICO:
		var elementoUnico ParametroElementoUnico
		if err := json.Unmarshal(d, &elementoUnico); err != nil {
			return err
		}
		p.Detalles = elementoUnico

	case TP_ELEMENTOS_COMPLETO:
		var elementosCompleto ParametroElementosCompleto
		if err := json.Unmarshal(d, &elementosCompleto); err != nil {
			return err
		}
		p.Detalles = elementosCompleto

	case TP_ELEMENTOS_PARCIAL:
		var elementos ParametroElementosParcial
		if err := json.Unmarshal(d, &elementos); err != nil {
			return err
		}
		p.Detalles = elementos
	}

	return nil
}
