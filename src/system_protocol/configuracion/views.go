package configuracion

import (
	"encoding/json"
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
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
	Nombre     string               `json:"nombre"`
	Template   string               `json:"bloqueTemplate"`
	Requisitos []string             `json:"requisitos"`
	Parametros map[string]Parametro `json:"informacion"`
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
	Condiciones      []CondicionTabla `json:"where"`
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

func CrearInfoViews(archivoJson string, bdd *b.Bdd, tablas []d.DescripcionTabla) (*v.InfoViews, error) {
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	var info InformacionViews
	if err := decodificador.Decode(&info); err != nil {
		return nil, fmt.Errorf("error al codificar tablas, con err: %v", err)
	}

	tablasPorNombre := make(map[string]*d.DescripcionTabla)
	for _, tabla := range tablas {
		tablasPorNombre[tabla.NombreTabla] = &tabla
	}

	views := make([]v.View, len(info.Views))
	for i, infoView := range info.Views {
		informaciones := make(map[string]v.Informacion)
		for clave := range infoView.Parametros {
			var informacion v.Informacion

			switch parametro := infoView.Parametros[clave].Parametro.(type) {
			case ParametroPathView:
				informacion = v.NewInformacionReferencia(parametro.View, map[string]string{})

			case ParametroElementos:
				if tabla, ok := tablasPorNombre[parametro.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", parametro.Tabla)

				} else {
					clavesRepresentativas := make(map[string]string)
					for _, condicion := range parametro.Condiciones {
						clavesRepresentativas[condicion.Clave] = condicion.Equal
					}
					referencias := make(map[string]v.InformacionReferencia)
					for _, referencia := range parametro.Referencias {
						referencias[referencia.Nombre] = v.NewInformacionReferencia(referencia.View, referencia.Requisitos)
					}

					// Hacer un chequeo con el template, o reduccion con el template, tal vez en vez de
					// tener que definirlo en el archivo, que este completamente obtenido por el template
					tiposPorClaves := make(map[string]d.TipoVariable)
					if len(parametro.ClavesSelectivas) == 0 {
						tiposPorClaves = tabla.TipoDadoClave
					} else {
						for _, clave := range parametro.ClavesSelectivas {
							if tipo, ok := tabla.TipoDadoClave[clave]; !ok {
								return nil, fmt.Errorf("la clave seleccionada %s no esta en la tabla", clave)
							} else {
								tiposPorClaves[clave] = tipo
							}
						}
					}

					parClaveRepresentacion := make(map[string]d.ElementoInformacion)
					for clave := range clavesRepresentativas {
						if tipo, ok := tabla.TipoDadoClave[clave]; !ok {
							return nil, fmt.Errorf("no hay tipo para la clave %s", clave)
						} else {
							parClaveRepresentacion[clave] = d.ElementoInformacion{
								Tipo:           tipo,
								Representacion: clavesRepresentativas[clave],
							}
						}
					}

					if query, err := d.NewQueryMultiples(tabla.NombreTabla, tiposPorClaves, parClaveRepresentacion, parametro.Ordenar); err != nil {
						return nil, err

					} else {
						informacion = v.NewInformacionTabla(query, referencias)
					}
				}

			case ParametroElementoUnico:
				if tabla, ok := tablasPorNombre[parametro.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", parametro.Tabla)

				} else {
					clavesRepresentativas := make(map[string]string)
					for _, condicion := range parametro.Condiciones {
						clavesRepresentativas[condicion.Clave] = condicion.Equal
					}

					// Hacer un chequeo con el template, o reduccion con el template, tal vez en vez de
					// tener que definirlo en el archivo, que este completamente obtenido por el template
					tiposPorClaves := make(map[string]d.TipoVariable)
					if len(parametro.ClavesSelectivas) == 0 {
						tiposPorClaves = tabla.TipoDadoClave
					} else {
						for _, clave := range parametro.ClavesSelectivas {
							if tipo, ok := tabla.TipoDadoClave[clave]; !ok {
								return nil, fmt.Errorf("la clave seleccionada %s no esta en la tabla", clave)
							} else {
								tiposPorClaves[clave] = tipo
							}
						}

					}

					parClaveRepresentacion := make(map[string]d.ElementoInformacion)
					for clave := range clavesRepresentativas {
						if tipo, ok := tabla.TipoDadoClave[clave]; !ok {
							return nil, fmt.Errorf("no hay tipo para la clave %s", clave)
						} else {
							parClaveRepresentacion[clave] = d.ElementoInformacion{
								Tipo:           tipo,
								Representacion: clavesRepresentativas[clave],
							}
						}
					}

					if query, err := d.NewQueryFila(tabla.NombreTabla, tiposPorClaves, parClaveRepresentacion, parametro.Ordenar); err != nil {
						return nil, err

					} else {
						informacion = v.NewInformacionFila(query)
					}
				}
			}

			informaciones[clave] = informacion
		}

		views[i] = v.NewView(bdd, infoView.Nombre, infoView.Template, infoView.Requisitos, informaciones)
	}

	return v.NewInfoViews(info.Inicio, views, info.PathTemplates, info.PathCss)
}
