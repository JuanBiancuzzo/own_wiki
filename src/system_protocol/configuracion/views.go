package configuracion

import (
	"encoding/json"
	"fmt"
	v "own_wiki/system_protocol/views"
	"strings"
)

type InformacionViews struct {
	Inicio        string `json:"inicio"`
	PathTemplates string `json:"templates"`
	PathCss       string `json:"css"`
	Views         []View `json:"views"`
}

type View struct {
	Nombre      string               `json:"nombre"`
	Template    InfoTemplete         `json:"template"`
	Requisitos  []string             `json:"requisitos"`
	Informacion map[string]Parametro `json:"informacion"`
}

type InfoTemplete struct {
	Archivo string `json:"archivo"`
	Bloque  string `json:"bloque"`
}

type TipoParametro string

const (
	TP_PATH_VIEW      = "pathView"
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

type ParametroElementos struct {
	HeaderParametro
	Tabla       string               `json:"tabla"`
	Condiciones []CondicionTabla     `json:"where"`
	Referencias []PathViewReferencia `json:"referencias"`
}

type ParametroElementoUnico struct {
	HeaderParametro
	Tabla       string           `json:"tabla"`
	Condiciones []CondicionTabla `json:"where"`
}

type CondicionTabla struct {
	Clave string `json:"clave"`
	Equal string `json:"equal"`
	From  string `json:"from,omitempty"`
}

type PathViewReferencia struct {
	Nombre     string            `json:"nombre"`
	View       string            `json:"view"`
	Requisitos map[string]string `json:"parametros"`
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

func CrearInfoViews(archivoJson string) (*v.InfoViews, error) {
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	var info InformacionViews
	if err := decodificador.Decode(&info); err != nil {
		return nil, fmt.Errorf("error al codificar tablas, con err: %v", err)
	}

	views := make([]v.View, len(info.Views))
	for i, infoView := range info.Views {

		views[i] = v.NewView(infoView.Nombre)
	}

	return v.NewInfoViews(info.Inicio, views, info.PathTemplates, info.PathCss)
}
