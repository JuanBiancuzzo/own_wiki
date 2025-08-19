package configuracion

import "encoding/json"

type TipoReferencia string

const (
	TR_PATH_VIEW = "pathView"
	TR_DELETE    = "eliminar"
	TR_UPDATE    = "actualizar"
)

type ReferenciaDeTabla struct {
	Referecnia any
}

type HeaderReferencia struct {
	Tipo TipoReferencia `json:"tipo"`
}

type PathViewReferencia struct {
	HeaderReferencia
	Nombre     string            `json:"nombre"`
	View       string            `json:"view"`
	Requisitos map[string]string `json:"parametros"`
}

// Se genera siempre para conseguir el id del elemento de la tabla
type UpdateReferencia struct {
	HeaderReferencia
	Nombre          string `json:"nombre"`
	BloqueRespuesta string `json:"bloque"`
}

// Se genera siempre para conseguir el id del elemento de la tabla
type DeleteReferencia struct {
	HeaderReferencia
	Nombre          string `json:"nombre"`
	BloqueRespuesta string `json:"bloque"`
}

func (r *ReferenciaDeTabla) UnmarshalJSON(d []byte) error {
	var header HeaderParametro
	if err := json.Unmarshal(d, &header); err != nil {
		return err
	}

	switch header.Tipo {
	case TR_PATH_VIEW:
		var referencia PathViewReferencia
		if err := json.Unmarshal(d, &referencia); err != nil {
			return err
		}
		r.Referecnia = referencia

	case TR_UPDATE:
		var referencia UpdateReferencia
		if err := json.Unmarshal(d, &referencia); err != nil {
			return err
		}
		r.Referecnia = referencia

	case TR_DELETE:
		var referencia DeleteReferencia
		if err := json.Unmarshal(d, &referencia); err != nil {
			return err
		}
		r.Referecnia = referencia
	}

	return nil
}
