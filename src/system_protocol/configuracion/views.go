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
	Nombre      string        `json:"nombre"`
	Template    string        `json:"bloqueTemplate"`
	Parametros  []string      `json:"parametros,omitempty"`
	Informacion []Informacion `json:"informacion,omitempty"`
}

func CrearInfoViews(archivoJson string, bdd *b.Bdd, tablas []d.DescripcionTabla) (*v.InfoViews, error) {
	decodificador := json.NewDecoder(strings.NewReader(archivoJson))
	pathView := v.NewPathView()

	var info InformacionViews
	if err := decodificador.Decode(&info); err != nil {
		return nil, fmt.Errorf("error al codificar tablas, con err: %v", err)
	}

	tablasPorNombre := make(map[string]*d.DescripcionTabla)
	for _, tabla := range tablas {
		tablasPorNombre[tabla.Nombre] = &tabla
	}

	endpoints := make(map[string]v.Endpoint)
	hayInicio := false

	for _, infoView := range info.Views {
		informaciones := make(map[string]v.FnInformacion)
		pathView.AgregarView(infoView.Nombre, infoView.Parametros)

		for _, variable := range infoView.Informacion {
			switch detalles := variable.Detalles.(type) {
			case ParametroElementoUnico:
				if informacion, err := crearInformacionElementoUnico(infoView.Parametros, detalles, tablasPorNombre); err != nil {
					return nil, err

				} else {
					informaciones[detalles.Nombre] = informacion
				}

			case ParametroElementosCompleto:
				tablas := make(InfoTablas)

				for _, infoTabla := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[infoTabla.Tabla]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", infoTabla.Tabla)

					} else if queryDato, err := infoTabla.CrearInformacionQuery(); err != nil {
						return nil, fmt.Errorf("hubo un error al obtener los detalles de la query, con error: %v", err)

					} else {
						tablas[tabla] = queryDato
					}
				}

				if informacion, err := crearInformacionElementosCompleto(tablas, infoView.Parametros, detalles.GroupBy, tablasPorNombre); err != nil {
					return nil, err

				} else {
					informaciones[detalles.Nombre] = informacion
				}

			case ParametroElementosParcial:
				tablas := make(InfoTablas)

				for _, infoTabla := range detalles.Tablas {
					if tabla, ok := tablasPorNombre[infoTabla.Tabla]; !ok {
						return nil, fmt.Errorf("no existe la tabla %s como una tabla registrada", infoTabla.Tabla)

					} else if queryDato, err := infoTabla.CrearInformacionQuery(); err != nil {
						return nil, fmt.Errorf("hubo un error al obtener los detalles de la query, con error: %v", err)

					} else {
						tablas[tabla] = queryDato
					}
				}

				if respuesta, err := crearInformacionElementosParcial(tablas, infoView.Parametros, detalles.GroupBy, tablasPorNombre); err != nil {
					return nil, err

				} else {
					informaciones[detalles.Nombre] = respuesta.Informacion
					for ruta := range respuesta.ExtraEndpoint {
						endpoints[ruta] = respuesta.ExtraEndpoint[ruta]
					}
				}
			}
		}

		ruta := infoView.Nombre
		endpoints[ruta] = v.NewView(bdd, infoView.Template, infoView.Parametros, informaciones)
		if ruta == info.Inicio {
			hayInicio = true
			endpoints[""] = v.NewView(bdd, infoView.Template, infoView.Parametros, informaciones)
		}

	}

	if !hayInicio {
		return nil, fmt.Errorf("no hay punto de inicio")
	}

	return v.NewInfoViews(endpoints, info.PathTemplates, info.PathCss, info.PathImagenes, pathView), nil
}

func crearInformacionElementoUnico(parametros []string, detalles ParametroElementoUnico, descripciones map[string]*d.DescripcionTabla) (v.FnInformacion, error) {
	tabla, ok := descripciones[detalles.Tabla]
	if !ok {
		return nil, fmt.Errorf("no existe la descripcion de la tabla %s", detalles.Tabla)
	}

	if !slices.Contains(parametros, detalles.PametroParaId) {
		return nil, fmt.Errorf("el id de la tabla no es uno de los parametros")
	}

	if query, err := d.NewQuerySimple(tabla, detalles.ClavesUsadas, detalles.PametroParaId, descripciones); err != nil {
		return nil, err

	} else {
		return v.NewInformacionFila(query, parametros)
	}
}

func crearInformacionElementosCompleto(tablas InfoTablas, parametros []string, groupBy []string, descripciones map[string]*d.DescripcionTabla) (v.FnInformacion, error) {
	if querys, err := d.NewQueryMultiples(tablas, groupBy, descripciones); err != nil {
		return nil, err

	} else {
		return v.NewInformacionCompleta(querys, parametros)
	}
}

func crearInformacionElementosParcial(tablas InfoTablas, parametros []string, groupBy []string, descripciones map[string]*d.DescripcionTabla) (v.RespuestaInformacion, error) {
	if querys, err := d.NewQueryMultiples(tablas, groupBy, descripciones); err != nil {
		return v.RespuestaInformacion{}, err

	} else {
		return v.NewInformacionParcial(querys, parametros)
	}
}
