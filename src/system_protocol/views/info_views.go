package views

import (
	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	PathTemplates string
	PathCss       string
	Endpoints     map[string]Endpoint
}

func NewInfoViews(endpoint map[string]Endpoint, pathTemplates, pathCss string) *InfoViews {
	return &InfoViews{
		Endpoints:     endpoint,
		PathTemplates: pathTemplates,
		PathCss:       pathCss,
	}
}

func (t *InfoViews) GenerarEndpoints(e *echo.Echo) {
	for ruta := range t.Endpoints {
		e.GET(ruta, echo.HandlerFunc(t.Endpoints[ruta]))
	}
}
