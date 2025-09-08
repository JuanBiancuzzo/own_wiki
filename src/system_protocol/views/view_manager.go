package views

import (
	"fmt"
	b "own_wiki/system_protocol/base_de_datos"
	"slices"

	"github.com/labstack/echo/v4"
)

type ViewManager struct {
	EndpointManager *EndpointManager
	HandlerView     map[string]echo.HandlerFunc

	viewLookUp []string
}

func NewViewManager() *ViewManager {
	return &ViewManager{
		EndpointManager: NewEndpointManager(),
		HandlerView:     make(map[string]echo.HandlerFunc),

		viewLookUp: []string{},
	}
}

func (vm *ViewManager) Agregar(view View) error {
	if slices.Contains(vm.viewLookUp, view.Nombre) {
		return fmt.Errorf("ya se cargo esa view")
	}

	if view.EsInicio {
		vm.HandlerView["/"] = echo.HandlerFunc(func(ec echo.Context) error {
			return ec.Render(200, fmt.Sprintf("%s/%s", view.Nombre, view.BloqueInicio), nil)
		})
	}

	vm.HandlerView[fmt.Sprintf("/%s", view.Nombre)] = echo.HandlerFunc(func(ec echo.Context) error {
		return ec.Render(200, fmt.Sprintf("%s/%s", view.Nombre, view.Bloque), nil)
	})

	vm.viewLookUp = append(vm.viewLookUp, view.Nombre)
	return view.RegistrarEndpoints(vm.EndpointManager)
}

func (vm *ViewManager) CreateURLPathView(view string) string {
	if !slices.Contains(vm.viewLookUp, view) {
		return "ERROR - No existe view"
	}

	return fmt.Sprintf("/%s", view)
}

func (vm *ViewManager) GenerarEndpoints(e *echo.Echo, bdd *b.Bdd) {
	for ruta := range vm.HandlerView {
		e.GET(ruta, vm.HandlerView[ruta])
	}

	vm.EndpointManager.GenerarEndpoints(e, bdd)
}
