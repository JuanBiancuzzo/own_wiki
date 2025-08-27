package configuracion

import (
	"encoding/json"
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"
	d "own_wiki/system_protocol/dependencias"
	v "own_wiki/system_protocol/views"
	"slices"
	"strings"
)

type InformacionViews struct {
	Inicio        string `json:"inicio"`
	PathTemplates string `json:"templates"`
	PathCss       string `json:"css"`
	PathImagenes  string `json:"imagenes"`
	Views         []View `json:"views"`
}

type InfoTablas map[*d.DescripcionTabla]d.InformacionQuery

type View struct {
	Nombre      string                 `json:"nombre"`
	Template    string                 `json:"bloqueTemplate"`
	Parametros  []string               `json:"parametros"`
	Informacion map[string]Informacion `json:"informacion"`
}

type RespuestaInformacionViews struct {
	InfoView *v.InfoViews
	PathView *v.PathView
}

func CrearInfoViews(archivoJson string, bdd *b.Bdd, tablas []d.DescripcionTabla) (*RespuestaInformacionViews, error) {
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))
	pathView := v.NewPathView()

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
		informaciones := make(map[string]v.FnInformacion)
		pathView.AgregarView(infoView.Nombre, infoView.Parametros)

		for nombreVariable := range infoView.Informacion {
			switch detalles := infoView.Informacion[nombreVariable].Detalles.(type) {
			case ParametroElementoUnico:
				if tabla, ok := tablasPorNombre[detalles.Tabla]; !ok {
					return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", detalles.Tabla)

				} else if informacion, err := crearInformacionElementoUnico(tabla, infoView.Parametros, detalles); err != nil {
					return nil, err

				} else {
					informaciones[nombreVariable] = informacion
				}

			case ParametroElementosCompleto:
				tablas := make(InfoTablas)

				for tablaUsada := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[tablaUsada]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", tablaUsada)

					} else if queryDato, err := detalles.Tablas[tablaUsada].CrearInformacionQuery(); err != nil {
						return nil, fmt.Errorf("hubo un error al obtener los detalles de la query, con error: %v", err)

					} else {
						tablas[tabla] = queryDato
					}
				}

				if informacion, err := crearInformacionElementosCompleto(tablas, infoView.Parametros, detalles.GroupBy); err != nil {
					return nil, err

				} else {
					informaciones[nombreVariable] = informacion
				}

			case ParametroElementosParcial:
				tablas := make(InfoTablas)

				for tablaUsada := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[tablaUsada]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", tablaUsada)

					} else if queryDato, err := detalles.Tablas[tablaUsada].CrearInformacionQuery(); err != nil {
						return nil, fmt.Errorf("hubo un error al obtener los detalles de la query, con error: %v", err)

					} else {
						tablas[tabla] = queryDato
					}
				}

				if respuesta, err := crearInformacionElementosParcial(tablas, infoView.Parametros, detalles.GroupBy); err != nil {
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
		PathView: pathView,
	}, nil
}

func crearInformacionElementoUnico(tabla *d.DescripcionTabla, parametros []string, detalles ParametroElementoUnico) (v.FnInformacion, error) {
	if !slices.Contains(parametros, detalles.PametroParaId) {
		return nil, fmt.Errorf("el id de la tabla no es uno de los parametros")
	}

	if query, err := d.NewQuerySimple(tabla, detalles.ClavesUsadas, detalles.PametroParaId); err != nil {
		return nil, err

	} else {
		return v.NewInformacionFila(query, parametros)
	}
}

func crearInformacionElementosCompleto(tablas InfoTablas, parametros []string, groupBy []string) (v.FnInformacion, error) {
	if querys, err := d.NewQueryMultiples(tablas, groupBy); err != nil {
		return nil, err

	} else {
		return v.NewInformacionCompleta(querys, parametros)
	}
}

func crearInformacionElementosParcial(tablas InfoTablas, parametros []string, groupBy []string) (v.RespuestaInformacion, error) {
	if querys, err := d.NewQueryMultiples(tablas, groupBy); err != nil {
		return v.RespuestaInformacion{}, err

	} else {
		return v.NewInformacionParcial(querys, parametros)
	}
}
