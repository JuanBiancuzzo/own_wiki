package views

import (
	"fmt"
)

type View struct {
	EsInicio     bool
	BloqueInicio string

	Nombre        string
	Bloque        string
	Endpoints     map[string]Endpoint
	PathTemplates []string
}

func NewView(esInicio bool, bloqueInicio, nombre, bloque string, endpoints map[string]Endpoint, templates []string) View {
	return View{
		EsInicio:     esInicio,
		BloqueInicio: bloqueInicio,

		Nombre:        nombre,
		Bloque:        bloque,
		Endpoints:     endpoints,
		PathTemplates: templates,
	}
}

func (v View) RegistrarEndpoints(endpointManager *EndpointManager) error {
	for rutaParcial := range v.Endpoints {
		ruta := fmt.Sprintf("%s/%s", v.Nombre, rutaParcial)
		if err := endpointManager.Agregar(ruta, v.Endpoints[rutaParcial]); err != nil {
			return err
		}
	}
	return nil
}
