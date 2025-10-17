package views

import (
	"fmt"
	"strings"

	b "github.com/JuanBiancuzzo/own_wiki/system_protocol/base_de_datos"

	"github.com/labstack/echo/v4"
)

type claves []string

type EndpointManager struct {
	Endpoints map[string]Endpoint

	clavesLookUp map[string]claves
}

func NewEndpointManager() *EndpointManager {
	return &EndpointManager{
		Endpoints:    make(map[string]Endpoint),
		clavesLookUp: make(map[string]claves),
	}
}

func (em *EndpointManager) Agregar(ruta string, endpoint Endpoint) error {
	if _, ok := em.clavesLookUp[ruta]; ok {
		return fmt.Errorf("ya se cargo ese parametro")
	}

	em.clavesLookUp[ruta] = endpoint.Parametros
	em.Endpoints[ruta] = endpoint
	return nil
}

func (em *EndpointManager) CreateURLPathEndpoint(endpoint string, valores ...any) string {
	if claves, ok := em.clavesLookUp[endpoint]; !ok {
		return fmt.Sprintf("ERROR - No existe endpoint '%s'", endpoint)

	} else if len(claves) != len(valores) {
		return fmt.Sprintf("ERROR - No suficientes parametros para '%s', deberian ser %d y se dieron %d", endpoint, len(claves), len(valores))

	} else {
		claveValor := make([]string, len(claves))
		for i, clave := range claves {
			valor := valores[i]
			claveValor[i] = fmt.Sprintf("%s=%v", clave, valor)
		}

		parametros := ""
		if len(claveValor) > 0 {
			parametros = fmt.Sprintf("?%s", strings.Join(claveValor, "&"))
		}

		return fmt.Sprintf("/%s%s", endpoint, parametros)
	}
}

func (em *EndpointManager) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	for ruta := range em.Endpoints {
		endpoint := em.Endpoints[ruta]

		handler := echo.HandlerFunc(endpoint.GenerarEndpoint(bdd))
		e.GET(ruta, handler)
	}
}
