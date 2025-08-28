package views

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	PathTemplates string
	PathCss       string
	PathImagenes  string
	PathView      *PathView
	Endpoints     map[string]Endpoint
}

func NewInfoViews(endpoint map[string]Endpoint, pathTemplates, pathCss, pathImagenes string, pathView *PathView) *InfoViews {
	return &InfoViews{
		Endpoints:     endpoint,
		PathTemplates: pathTemplates,
		PathCss:       pathCss,
		PathImagenes:  pathImagenes,
		PathView:      pathView,
	}
}

func (t *InfoViews) GenerarEndpoints(e *echo.Echo) {
	for ruta := range t.Endpoints {
		e.GET(fmt.Sprintf("/%s", ruta), echo.HandlerFunc(t.Endpoints[ruta]))
	}
}
