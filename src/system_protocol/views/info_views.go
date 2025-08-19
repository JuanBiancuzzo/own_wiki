package views

import (
	"fmt"

	"github.com/labstack/echo/v4"
)

type InfoViews struct {
	Inicio        int
	PathTemplates string
	PathCss       string
	Views         []View
}

func NewInfoViews(inicio string, views []View, pathTemplates, pathCss string) (*InfoViews, error) {
	for i, view := range views {
		if view.Nombre == inicio {
			return &InfoViews{
				Inicio:        i,
				Views:         views,
				PathTemplates: pathTemplates,
				PathCss:       pathCss,
			}, nil
		}
	}

	return nil, fmt.Errorf("no se establecio un endpoint inicial")
}

func (t *InfoViews) GenerarEndpoints(e *echo.Echo) {
	for i, view := range t.Views {
		ruta := fmt.Sprintf("/%s", view.Nombre)
		if i == t.Inicio {
			ruta = "/"
		}

		view.GenerarEndpoint(ruta, e)
	}
}
