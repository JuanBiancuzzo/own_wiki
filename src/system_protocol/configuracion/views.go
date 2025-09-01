package configuracion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	d "own_wiki/system_protocol/dependencias"
	v "own_wiki/system_protocol/views"
	"slices"
	"strings"
)

type InformacionViews struct {
	Inicio             string     `json:"inicio"`
	PathCss            string     `json:"css"`
	PathImagenes       string     `json:"imagenes"`
	PathViews          []string   `json:"views"`
	EndpointsGenerales []Endpoint `json:"endpointsGenerales"`
}

type View struct {
	Nombre         string     `json:"nombre"`
	Templates      []string   `json:"templates"`
	BloqueTemplate string     `json:"bloqueTemplate"`
	Endpoints      []Endpoint `json:"endpoints"`

	esInicio bool
	pathView string
}

type InfoTablas map[*d.DescripcionTabla]d.InformacionQuery

type Endpoint struct {
	Nombre         string        `json:"nombre"`
	BloqueTemplate string        `json:"bloqueTemplate"`
	Parametros     []string      `json:"parametros,omitempty"`
	Informacion    []Informacion `json:"informacion,omitempty"`
}

func leerView(pathView string) (View, error) {
	if bytesView, err := os.ReadFile(pathView); err != nil {
		return View{}, fmt.Errorf("error al leer el archivo de configuracion para las views, con error: %v", err)

	} else {
		var view View
		carpetas := strings.Split(pathView, "/")
		view.pathView = strings.Join(carpetas[:len(carpetas)-1], "/")

		decodificador := json.NewDecoder(bytes.NewReader(bytesView))

		if err := decodificador.Decode(&view); err != nil {
			return view, fmt.Errorf("error al codificar tablas, con err: %v", err)
		}

		return view, nil
	}
}

func CrearInfoViews(pathConfiguracion string, tablas []d.DescripcionTabla) (*v.InfoViews, error) {
	var informacionViews InformacionViews
	var err error

	if bytesJson, err := os.ReadFile(fmt.Sprintf("%s/%s", pathConfiguracion, "views.json")); err != nil {
		return nil, fmt.Errorf("error al leer el archivo de configuracion para las views, con error: %v", err)

	} else {
		decodificador := json.NewDecoder(bytes.NewReader(bytesJson))
		if err := decodificador.Decode(&informacionViews); err != nil {
			return nil, fmt.Errorf("error al codificar tablas, con err: %v", err)
		}
	}

	cantidadViews := len(informacionViews.PathViews)
	if cantidadViews == 0 {
		return nil, fmt.Errorf("no se ingresaron views para el proyecto, recordar poner el array de 'views'")
	}

	hayInicio := false
	viewsInfo := make([]View, cantidadViews)
	for i := range cantidadViews {
		if viewsInfo[i], err = leerView(fmt.Sprintf("%s/%s.json", pathConfiguracion, informacionViews.PathViews[i])); err != nil {
			return nil, err
		}

		if viewsInfo[i].Nombre == informacionViews.Inicio {
			hayInicio = true
			viewsInfo[i].esInicio = true
		}
	}

	if !hayInicio {
		return nil, fmt.Errorf("no se ingresó una view la cual corresponda ser el inicio")
	}

	tablasPorNombre := make(map[string]*d.DescripcionTabla)
	for _, tabla := range tablas {
		tablasPorNombre[tabla.Nombre] = &tabla
	}

	views := make([]v.View, len(viewsInfo))
	for i, viewInfo := range viewsInfo {
		endpoints := make(map[string]v.Endpoint)

		for _, infoEndpoint := range viewInfo.Endpoints {
			informaciones := make(map[string]v.FnInformacion)

			for _, variable := range infoEndpoint.Informacion {
				switch detalles := variable.Detalles.(type) {
				case ParametroElementoUnico:
					if informacion, err := crearInformacionElementoUnico(infoEndpoint.Parametros, detalles, tablasPorNombre); err != nil {
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

					if informacion, err := crearInformacionElementosCompleto(tablas, infoEndpoint.Parametros, detalles.GroupBy, tablasPorNombre); err != nil {
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

					if respuesta, err := crearInformacionElementosParcial(tablas, infoEndpoint.Parametros, detalles.GroupBy, tablasPorNombre); err != nil {
						return nil, err

					} else {
						informaciones[detalles.Nombre] = respuesta.Informacion
						for ruta := range respuesta.ExtraEndpoint {
							endpoints[ruta] = respuesta.ExtraEndpoint[ruta]
						}
					}
				}
			}

			endpoints[infoEndpoint.Nombre] = v.NewEndpoint(infoEndpoint.BloqueTemplate, infoEndpoint.Parametros, informaciones)
		}

		pathsTemplate := make([]string, len(viewInfo.Templates))
		for i, pathTemplate := range viewInfo.Templates {
			pathsTemplate[i] = fmt.Sprintf("%s/%s", viewInfo.pathView, pathTemplate)
		}
		views[i] = v.NewView(viewInfo.esInicio, viewInfo.Nombre, viewInfo.BloqueTemplate, endpoints, pathsTemplate)
	}

	return v.NewInfoViews(views, informacionViews.PathCss, informacionViews.PathImagenes)
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
