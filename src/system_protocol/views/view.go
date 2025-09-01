package views

import (
	"fmt"
	b "own_wiki/system_protocol/bass_de_datos"

	"github.com/labstack/echo/v4"
)

type View struct {
	Nombre        string
	Bloque        string
	Endpoints     map[string]Endpoint
	PathTemplates []string
}

func NewView(nombre, bloque string, endpoints map[string]Endpoint, templates []string) View {
	return View{
		Nombre:        nombre,
		Bloque:        bloque,
		Endpoints:     endpoints,
		PathTemplates: templates,
	}
}

func (v View) RegistrarEndpoints(pathView *PathView) {
	for ruta := range v.Endpoints {
		pathView.AgregarView(ruta, v.Endpoints[ruta].Parametros)
	}
}

func (v View) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	handler := echo.HandlerFunc(func(ec echo.Context) error {
		return ec.Render(200, v.Bloque, nil)
	})
	e.GET(fmt.Sprintf("/%s", v.Nombre), handler)

	for ruta := range v.Endpoints {
		handler = echo.HandlerFunc(v.Endpoints[ruta].GenerarEndpoint(bdd, v.Nombre))
		e.GET(fmt.Sprintf("/%s", v.Nombre), handler)
	}
}
