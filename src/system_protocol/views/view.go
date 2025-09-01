package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type View struct {
	EsInicio      bool
	Nombre        string
	Bloque        string
	Endpoints     map[string]Endpoint
	PathTemplates []string
}

func NewView(esInicio bool, nombre, bloque string, endpoints map[string]Endpoint, templates []string) View {
	return View{
		EsInicio:      esInicio,
		Nombre:        nombre,
		Bloque:        bloque,
		Endpoints:     endpoints,
		PathTemplates: templates,
	}
}

func (v View) RegistrarEndpoints(pathView *PathEndpoint) error {
	for ruta := range v.Endpoints {
		if err := pathView.AgregarEndpoint(ruta, v.Endpoints[ruta].Parametros); err != nil {
			return err
		}
	}
	return nil
}

func (v View) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	handler := echo.HandlerFunc(func(ec echo.Context) error {
		return ec.Render(200, fmt.Sprintf("%s/%s", v.Nombre, v.Bloque), nil)
	})
	e.GET(fmt.Sprintf("/%s", v.Nombre), handler)
	if v.EsInicio {
		e.GET("/", handler)
	}

	for ruta := range v.Endpoints {
		handler = echo.HandlerFunc(v.Endpoints[ruta].GenerarEndpoint(bdd, v.Nombre))
		e.GET(fmt.Sprintf("/%s/%s", v.Nombre, ruta), handler)
	}
}
