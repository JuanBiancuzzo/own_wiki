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

type InfoTablas map[*d.DescripcionTabla]InformacionTabla

type View struct {
	Nombre      string                 `json:"nombre"`
	Template    string                 `json:"bloqueTemplate"`
	Parametros  []string               `json:"parametros"`
	Informacion map[string]Informacion `json:"informacion"`
}

type RespuestaInformacion struct {
	Informacion   v.Informacion
	ExtraEndpoint map[string]v.Endpoint
}

type RespuestaInformacionViews struct {
	InfoView *v.InfoViews
	PathView *v.PathView
}

func CrearInfoViews(archivoJson string, bdd *b.Bdd, tablas []d.DescripcionTabla) (*RespuestaInformacionViews, error) {
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))

	var info InformacionViews
	if err := decodificador.Decode(&info); err != nil {
		return nil, fmt.Errorf("error al codificar tablas, con err: %v", err)
	}

	tablasPorNombre := make(map[string]*d.DescripcionTabla)
	for _, tabla := range tablas {
		tablasPorNombre[tabla.NombreTabla] = &tabla
	}

	endpoints := make(map[string]v.Endpoint)
	hayInicio := false

	for _, infoView := range info.Views {
		informaciones := make(map[string]v.Informacion)
		for nombreVariable := range infoView.Informacion {
			switch detalles := infoView.Informacion[nombreVariable].Detalles.(type) {
			case ParametroElementoUnico:
				if tabla, ok := tablasPorNombre[detalles.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", detalles.Tabla)

				} else if informacion, err := crearInformacionElementoUnico(tabla, infoView, nombreVariable, detalles); err != nil {
					return nil, err

				} else {
					informaciones[nombreVariable] = informacion
				}

			case ParametroElementosCompleto:
				tablas := make(map[*d.DescripcionTabla]InformacionTabla)
				for tablaUsada := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[tablaUsada]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", tablaUsada)

					} else {
						tablas[tabla] = detalles.Tablas[tablaUsada]
					}
				}

				if respuesta, err := crearInformacionElementosCompleto(infoView, tablas, nombreVariable, detalles); err != nil {
					return nil, err

				} else {
					informaciones[nombreVariable] = respuesta.Informacion
					for ruta := range respuesta.ExtraEndpoint {
						endpoints[ruta] = respuesta.ExtraEndpoint[ruta]
					}
				}

			case ParametroElementosParcial:
				tablas := make(map[*d.DescripcionTabla]InformacionTabla)
				for tablaUsada := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[tablaUsada]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", tablaUsada)

					} else {
						tablas[tabla] = detalles.Tablas[tablaUsada]
					}
				}

				if respuesta, err := crearInformacionElementosParcial(infoView, tablas, nombreVariable, detalles); err != nil {
					return nil, err

				} else {
					informaciones[nombreVariable] = respuesta.Informacion
					for ruta := range respuesta.ExtraEndpoint {
						endpoints[ruta] = respuesta.ExtraEndpoint[ruta]
					}
				}
			}
		}

		ruta := infoView.Nombre
		if ruta == info.Inicio {
			hayInicio = true
			ruta = ""
		}

		endpoints[ruta] = v.NewView(bdd, infoView.Template, infoView.Parametros, informaciones)
	}

	if !hayInicio {
		return nil, fmt.Errorf("no hay punto de inicio")
	}

	return &RespuestaInformacionViews{
		InfoView: v.NewInfoViews(endpoints, info.PathTemplates, info.PathCss),
		PathView: v.NewPathView(),
	}, nil
}

func crearInformacionElementosCompleto(infoView View, tablas InfoTablas, nombreVariable string, parametro ParametroElementosCompleto) (RespuestaInformacion, error) {

	return RespuestaInformacion{
		Informacion:   New,
		ExtraEndpoint: make(map[string]v.Endpoint),
	}, nil
}

func crearInformacionElementosParcial(infoView View, tablas InfoTablas, nombreVariable string, parametro ParametroElementosParcial) (RespuestaInformacion, error) {

	return RespuestaInformacion{
		ExtraEndpoint: make(map[string]v.Endpoint),
	}, nil
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
