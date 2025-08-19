package configuracion

import "encoding/json"

type TipoParametro string

const (
	TP_PATH_VIEW      = "pathView"
	TP_INSERTAR       = "insertar"
	TP_ELEMENTOS      = "elementosTabla"
	TP_ELEMENTO_UNICO = "elementoUnicoTabla"
)

type Parametro struct {
	Parametro any
}

type HeaderParametro struct {
	Tipo TipoParametro `json:"tipo"`
}

type ParametroPathView struct {
	HeaderParametro
	View string `json:"view"`
}

type ParametroInsertar struct {
	HeaderParametro
	Tabla           string            `json:"table"`
	Parametros      map[string]string `json:"parametros"`
	BloqueRespuesta string            `json:"bloque"`
}

type ParametroElementos struct {
	HeaderParametro
	Tabla            string               `json:"tabla"`
	Condiciones      []CondicionTabla     `json:"condicion"`
	Ordenar          []string             `json:"ordenar"`
	ClavesSelectivas []string             `json:"claves"`
	Referencias      []PathViewReferencia `json:"referencias"`
}

type ParametroElementoUnico struct {
	HeaderParametro
	Tabla            string           `json:"tabla"`
	Ordenar          []string         `json:"ordenar"`
	ClavesSelectivas []string         `json:"claves"`
	Condiciones      []CondicionTabla `json:"condicion"`
}

type CondicionTabla struct {
	Clave string `json:"clave"`
	Equal string `json:"equal"`
}

func (p *Parametro) UnmarshalJSON(d []byte) error {
	var header HeaderParametro
	if err := json.Unmarshal(d, &header); err != nil {
		return err
	}

	switch header.Tipo {
	case TP_PATH_VIEW:
		var pathView ParametroPathView
		if err := json.Unmarshal(d, &pathView); err != nil {
			return err
		}
		p.Parametro = pathView

	case TP_INSERTAR:
		var insertar ParametroInsertar
		if err := json.Unmarshal(d, &insertar); err != nil {
			return err
		}
		p.Parametro = insertar

	case TP_ELEMENTOS:
		var elementos ParametroElementos
		if err := json.Unmarshal(d, &elementos); err != nil {
			return err
		}
		p.Parametro = elementos

	case TP_ELEMENTO_UNICO:
		var elementoUnico ParametroElementoUnico
		if err := json.Unmarshal(d, &elementoUnico); err != nil {
			return err
		}
		p.Parametro = elementoUnico
	}

	return nil
}
