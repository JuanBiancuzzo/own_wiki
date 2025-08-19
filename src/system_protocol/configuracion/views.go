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
			switch parametro := infoView.Parametros[clave].Parametro.(type) {
			case ParametroPathView:
				informaciones[clave] = v.NewInformacionReferencia(parametro.View, map[string]string{})

			case ParametroElementos:
				if tabla, ok := tablasPorNombre[parametro.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", parametro.Tabla)

				} else if informacion, err := crearInformacionElementos(tabla, infoView, clave, parametro); err != nil {
					return nil, err

				} else {
					informaciones[clave] = informacion
				}

			case ParametroElementoUnico:
				if tabla, ok := tablasPorNombre[parametro.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", parametro.Tabla)

				} else if informacion, err := crearInformacionElementoUnico(tabla, infoView, clave, parametro); err != nil {
					return nil, err

				} else {
					informaciones[clave] = informacion
				}
			}
		}

		views[i] = v.NewView(bdd, infoView.Nombre, infoView.Template, infoView.Requisitos, informaciones, map[string]v.Endpoint{})
	}

	return v.NewInfoViews(info.Inicio, views, info.PathTemplates, info.PathCss)
}

func crearInformacionElementos(tabla *d.DescripcionTabla, infoView View, clave string, parametro ParametroElementos) (v.Informacion, error) {
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
		return nil, fmt.Errorf("en la view %s, en la clave: %s se tuvo: %v", infoView.Nombre, clave, err)

	} else {
		return v.NewInformacionTabla(clave, query, referencias), nil
	}
}

func crearInformacionElementoUnico(tabla *d.DescripcionTabla, infoView View, clave string, parametro ParametroElementoUnico) (v.Informacion, error) {
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
		return nil, fmt.Errorf("en la view %s, en la clave: %s se tuvo: %v", infoView.Nombre, clave, err)

	} else {
		return v.NewInformacionFila(query), nil
	}
}
